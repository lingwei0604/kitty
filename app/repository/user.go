package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"git.yingzhongshare.com/mkt/kitty/app/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	fr *FileRepo
	db *gorm.DB
}

var ErrAlreadyBind = errors.New("third party account is bound to another user")
var ErrRecordNotFound = errors.New("record not found")

func IsErrRecordNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}

const emsg = "UserRepo"

func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return errors.Wrap(err, emsg)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, id uint, user entity.User) (newUser *entity.User, err error) {
	var u entity.User

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var tmp entity.User
		tx = tx.WithContext(ctx)
		if err := tx.First(&u, "id = ?", id).Error; err != nil {
			return errors.Wrap(err, fmt.Sprintf("while find user by id: %d", id))
		}
		if user.Mobile.Valid {
			err := tx.Where("package_name = ? and mobile = ?", u.PackageName, user.Mobile.String).First(&tmp).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "while looking for mobile name")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) && tmp.ID != id {
				tmp.Mobile = sql.NullString{Valid: false}
				if err := tx.Save(&tmp).Error; err != nil {
					return errors.Wrap(err, "unable to remove mobile from other users")
				}
			}
		}
		if user.WechatOpenId.Valid {
			err := tx.Where("package_name = ? and wechat_open_id = ?", u.PackageName, user.WechatOpenId.String).First(&tmp).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "while looking for wechat open id")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) && tmp.ID != id {
				tmp.WechatOpenId = sql.NullString{Valid: false}
				tmp.WechatUnionId = sql.NullString{Valid: false}
				tmp.WechatExtra = nil
				if err := tx.Save(&tmp).Error; err != nil {
					return errors.Wrap(err, "unable to remove wechat from other users")
				}
			}
		}

		if user.TaobaoOpenId.Valid {
			err := tx.Where("package_name = ? and taobao_open_id = ?", u.PackageName, user.TaobaoOpenId.String).First(&tmp).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "while looking for taobao open id")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) && tmp.ID != id {
				tmp.TaobaoOpenId = sql.NullString{Valid: false}
				tmp.TaobaoExtra = nil
				if err := tx.Save(&tmp).Error; err != nil {
					return errors.Wrap(err, "unable to remove taobao from other users")
				}
			}
		}
		err := tx.Model(entity.User{}).Unscoped().Where("id = ?", id).Updates(user).Error
		if err != nil {
			if IsDuplicateErr(err) {
				return ErrAlreadyBind
			}
			return err
		}

		if err := tx.First(&u, "id = ?", id).Error; err != nil {
			return errors.Wrap(err, fmt.Sprintf("while find user by id: %d", id))
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction rolled back")
	}
	return &u, nil
}

func (r *UserRepo) Delete(ctx context.Context, id uint) (err error) {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *UserRepo) Exists(ctx context.Context, id uint) bool {
	var (
		u entity.User
	)
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error

	if err != nil {
		return false
	}

	return true
}

func (r *UserRepo) UpdateCallback(ctx context.Context, id uint, f func(user *entity.User) error) (err error) {
	var u entity.User
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		err := tx.Model(entity.User{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&u).Error
		if err != nil {
			return errors.Wrap(err, emsg)
		}
		err = f(&u)
		if err != nil {
			return err
		}

		err = tx.Save(u).Error
		if err != nil {
			if IsDuplicateErr(err) {
				return ErrAlreadyBind
			}
			return errors.Wrap(err, emsg)
		}
		return nil
	})
}

func NewUserRepo(db *gorm.DB, fr *FileRepo, assigner entity.IDAssigner) *UserRepo {
	return &UserRepo{fr, db.Set("IDAssigner", assigner)}
}

func (r *UserRepo) GetUserByOpenID(ctx context.Context, packageName string, openID string) (*entity.User, error) {
	var (
		u entity.User
	)
	err := r.db.WithContext(ctx).First(&u, "package_name = ? and wechat_open_id = ?", packageName, openID).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}

func (r *UserRepo) GetFromWechat(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) (*entity.User, error) {
	var (
		u entity.User
	)

	wechatUser.CommonSUUID = device.Suuid
	wechatUser.CommonSMID = device.SMID
	wechatUser.PackageName = packageName
	wechatUser.WechatOpenId = sql.NullString{String: wechat, Valid: true}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			u.AddNewDevice(device)
			tx.WithContext(ctx).Save(device)
		}()
		err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"package_name = ? and wechat_open_id = ?", packageName, wechat,
		).First(&u).Error

		if err == nil {
			if u.CommonSUUID != device.Suuid {
				u.CommonSUUID = device.Suuid
				u.CommonSMID = device.SMID
				err = tx.Save(u).Error
				if err != nil {
					return errors.Wrap(err, "GetFromWechat: unable to update user with new suuid")
				}
			}
			return nil
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}
		if wechatUser.HeadImg != "" {
			wechatUser.HeadImg, _ = r.fr.UploadFromUrl(ctx, wechatUser.HeadImg)
		}
		stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: buildAssignmentColumns(wechatUser),
		}).Create(&wechatUser)
		if err := stmt.Error; err != nil {
			return err
		}
		wechatUser.IsNew = stmt.RowsAffected == 1
		u = wechatUser
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}

func (r *UserRepo) GetFromMobile(ctx context.Context, packageName, mobile string, device *entity.Device) (*entity.User, error) {
	var (
		mobileUser entity.User
		u          entity.User
	)

	mobileUser.CommonSUUID = device.Suuid
	mobileUser.CommonSMID = device.SMID
	mobileUser.PackageName = packageName
	mobileUser.Mobile = sql.NullString{String: mobile, Valid: true}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			u.AddNewDevice(device)
			tx.WithContext(ctx).Save(device)
		}()
		err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"package_name = ? and mobile = ?", packageName, mobile,
		).First(&u).Error

		if err == nil {
			if u.CommonSUUID != device.Suuid {
				u.CommonSUUID = device.Suuid
				u.CommonSMID = device.SMID
				err = tx.Save(u).Error
				if err != nil {
					return errors.Wrap(err, "GetFromMobile: unable to update user with new suuid")
				}
			}
			return nil
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}
		stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: buildAssignmentColumns(mobileUser),
		}).Create(&mobileUser)
		if err := stmt.Error; err != nil {
			return err
		}
		mobileUser.IsNew = stmt.RowsAffected == 1
		u = mobileUser
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}

func (r *UserRepo) GetRecentUserByDevice(ctx context.Context, packageName string, device *entity.Device) (*entity.User, error) {
	rr := &deviceRepo{r.db}
	return rr.GetRecentUserByDevice(ctx, packageName, device)
}

func (r *UserRepo) GetFromDevice(ctx context.Context, packageName, suuid string, device *entity.Device) (*entity.User, error) {
	var (
		err error
		u   *entity.User
	)
	err = r.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			u.AddNewDevice(device)
			tx.WithContext(ctx).Save(device)
		}()

		// 先使用suuid进行查询
		err = tx.WithContext(ctx).Where("package_name = ? and common_s_uuid = ?", packageName, suuid).First(&u).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if !u.Empty() {
			return nil
		}
		// 再尝试设备
		rr := &deviceRepo{tx}
		u, err = rr.GetRecentUserByDevice(ctx, packageName, &entity.Device{Suuid: suuid, Oaid: device.Oaid, Idfa: device.Idfa, AndroidId: device.AndroidId})
		if err != nil && !errors.Is(err, ErrRecordNotFound) {
			return err
		}
		if !u.Empty() {
			return nil
		}
		// 否则创建新用户
		u = &entity.User{PackageName: packageName, CommonSUUID: suuid, CommonSMID: device.SMID}
		err = tx.WithContext(ctx).Create(&u).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return u, nil
}

func (r *UserRepo) Get(ctx context.Context, id uint) (*entity.User, error) {
	var (
		u entity.User
	)
	if err := r.db.WithContext(ctx).Where("id = ?", id).Find(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	if u.ID == 0 {
		return nil, ErrRecordNotFound
	}
	return &u, nil
}

func (r *UserRepo) DeleteDevices(ctx context.Context, uid uint) ([]entity.Device, error) {
	var d []entity.Device
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("user_id = ?", uid).Delete(&d).Error; err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return d, nil
}

func (r *UserRepo) DeleteDevicesByOaid(ctx context.Context, oaid string) ([]entity.Device, error) {
	return r.deleteDevicesBy(ctx, "oaid", oaid)
}
func (r *UserRepo) DeleteDevicesByIdfa(ctx context.Context, idfa string) ([]entity.Device, error) {
	return r.deleteDevicesBy(ctx, "idfa", idfa)
}
func (r *UserRepo) DeleteDevicesByAndroid(ctx context.Context, android string) ([]entity.Device, error) {
	return r.deleteDevicesBy(ctx, "android_id", android)
}
func (r *UserRepo) deleteDevicesBy(ctx context.Context, key, value string) ([]entity.Device, error) {
	if key != "suuid" && key != "oaid" && key != "idfa" && key != "android_id" {
		return nil, fmt.Errorf("不允许的字段")
	}
	if value == "" {
		return nil, fmt.Errorf("不允许的值")
	}
	var d []entity.Device
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where(fmt.Sprintf("%s = ?", key), value).Delete(&d).Error; err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return d, nil
}

func (r *UserRepo) GetAll(ctx context.Context, where ...clause.Expression) ([]entity.User, error) {
	var (
		u []entity.User
	)
	if err := r.db.WithContext(ctx).Clauses(where...).Find(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	return u, nil
}

func (r *UserRepo) GetIDAndCreatedAtBySUUID(ctx context.Context, packageName, suuid string) (uint, time.Time, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).Select("id", "created_at").Where("package_name = ? and common_s_uuid = ?", packageName, suuid).Order("created_at desc").First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, time.Time{}, ErrRecordNotFound
		}
		return 0, time.Time{}, errors.Wrap(err, emsg)
	}
	return u.ID, u.CreatedAt, nil
}

func (r *UserRepo) GetByDevice(ctx context.Context, packageName string, device entity.Device) ([]entity.User, error) {
	var (
		u []entity.User
	)
	statement := r.db.WithContext(ctx).Table("kitty_users").Select("DISTINCT kitty_users.*").Joins("LEFT JOIN kitty_devices ON kitty_devices.user_id = kitty_users.id").Where("kitty_users.package_name = ?", packageName)
	if device.Imei != "" {
		statement = statement.Where("kitty_devices.imei = ?", device.Imei)
	}
	if device.Oaid != "" {
		statement = statement.Where("kitty_devices.oaid = ?", device.Oaid)
	}
	if device.Idfa != "" {
		statement = statement.Where("kitty_devices.idfa = ?", device.Idfa)
	}
	if device.Suuid != "" {
		statement = statement.Where("kitty_users.common_s_uuid = ?", device.Suuid)
	}

	if err := statement.Order("id desc").Find(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	return u, nil
}

func (r *UserRepo) Count(ctx context.Context, where ...clause.Expression) (int64, error) {
	if len(where) == 0 {
		return r.CountAll(ctx)
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Clauses(where...).Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, ErrRecordNotFound
		}
		return 0, errors.Wrap(err, emsg)
	}
	return count, nil
}

func (r *UserRepo) CountAll(ctx context.Context) (int64, error) {
	var count int64
	handle, _ := r.db.DB()
	rows, err := handle.QueryContext(
		ctx,
		`select table_rows from information_schema.TABLES where TABLE_NAME="kitty_users" and table_schema="monetization"`,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepo) GetByWechat(ctx context.Context, packageName, wechat string) (user *entity.User, err error) {
	var (
		u entity.User
	)
	if err := r.db.WithContext(ctx).Where(
		"package_name = ? and wechat_open_id = ?", packageName, wechat,
	).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}

	return &u, err
}

func buildAssignmentColumns(user entity.User) clause.Set {
	var columns []string
	if user.Mobile.Valid {
		columns = append(columns, "mobile")
	}
	if user.WechatOpenId.Valid {
		columns = append(columns, "wechat_open_id")
	}
	if user.WechatUnionId.Valid {
		columns = append(columns, "wechat_union_id")
	}
	if user.VersionCode != "" {
		columns = append(columns, "version_code")
	}
	if user.UserName != "" {
		columns = append(columns, "user_name")
	}
	if user.HeadImg != "" {
		columns = append(columns, "head_img")
	}
	if user.Gender != 0 {
		columns = append(columns, "gender")
	}
	if user.CommonSMID != "" {
		columns = append(columns, "common_sm_id")
	}
	if user.Birthday != "" {
		columns = append(columns, "birthday")
	}
	if user.TaobaoOpenId.Valid {
		columns = append(columns, "taobao_open_id")
	}

	return clause.AssignmentColumns(columns)
}

type deviceRepo struct {
	db *gorm.DB
}

func (r *deviceRepo) GetRecentUserByDevice(ctx context.Context, packageName string, device *entity.Device) (*entity.User, error) {
	// 先通过 suuid/[oaid,android_id]/idfa 查找最新的用户id

	u, err := r.getByDeviceWithoutJoin(ctx, packageName, device)
	if err != nil {
		return nil, err
	}

	if u.CommonSUUID != device.Suuid && device.Suuid != "" {
		u.CommonSUUID = device.Suuid
		if u.CommonSUUID != "" {
			r.db.WithContext(ctx).Model(&u).Update("common_s_uuid", u.CommonSUUID)
		}
	}
	return u, nil
}

// getByDeviceWithoutJoin 查询oaid对应的用户，不进行表关联
func (r *deviceRepo) getByDeviceWithoutJoin(ctx context.Context, packageName string, device *entity.Device) (*entity.User, error) {
	var (
		suuid = device.Suuid
		// android
		oaid      = device.GetOaid()
		androidId = device.GetAndroidId()
		// ios
		idfa = device.GetIdfa()
	)

	if suuid == "" && oaid == "" && androidId == "" && idfa == "" {
		return nil, ErrRecordNotFound
	}
	var (
		u   entity.User
		uid uint
	)

	// 优先以suuid查询
	statement := r.db.WithContext(ctx).Table("kitty_devices")
	orState := r.db.WithContext(ctx)
	if suuid != "" {
		orState = orState.Or("suuid = ?", suuid)
	}
	if androidId != "" {
		orState = orState.Or("android_id = ?", androidId)
	}
	if oaid != "" {
		orState = orState.Or("oaid = ?", oaid)
	}
	// android 与 ios本身互斥
	if idfa != "" {
		orState = orState.Or("idfa = ?", idfa)
	}

	if err := statement.Where(r.db.WithContext(ctx).Where("deleted_at is null").Where(orState)).
		Select("ifnull(max(user_id),0)").Find(&uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}

	if uid == 0 {
		return nil, ErrRecordNotFound
	}

	if err := r.db.WithContext(ctx).Table("kitty_users").Where("id = ? and package_name=?", uid, packageName).Order("id desc").Find(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}
