import {check, group, sleep} from 'k6';
import http from 'k6/http';
import {JWT} from "./jwt.js";

function getRandomInt(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min)) + min; //不含最大值，含最小值
}

export let options = {
  vus: 500,
  duration: '15s',
  thresholds: {
    checks: ['rate>0.999'],
    http_req_duration: ['p(90)<200'], // 95% of requests should be below 200ms
  },
};

export default function () {
  let obtainedJWT, userId
  let suuid = getRandomInt(1, 99999).toString()

  group("app tests", function () {
    // 模拟登陆
    {
      let url = `http://${__ENV.MY_HOSTNAME}/app/v2/login`;
      let payload = JSON.stringify({
        "device": {
          "imei": "string",
          "idfa": "string",
          "android_id": "string",
          "suuid": suuid,
          "mac": "string",
          "os": "OS_UNKNOWN",
          "oaid": "string",
          "smid": "string"
        },
        "channel": "performance testing",
        "version_code": "1000",
        "package_name": "com.performance.testing",
        "third_party_id": "string"
      });

      let params = {
        headers: {
          'Content-Type': 'application/json',
        },
      };

      let res = http.post(url, payload, params);
      let data = res.json()
      check(res, {
        'is status 200': (r) => r.status === 200,
        'code is 0': (r) => data.code === 0,
        'jwt obtained': () => data.data && data.data.token
      });
      if (data != null && data.data != null) {
        obtainedJWT = data.data.token
        userId = data.data.id
      }
      if (!obtainedJWT) {
        console.log(res.body)
        return
      }
    }

    // 模拟请求用户信息
    {
      sleep(Math.random() * 3)
      let url = `http://${__ENV.MY_HOSTNAME}/app/v2/info/0`
      let params = {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': "bearer " + obtainedJWT
        },
      };
      let res = http.get(url, params);
      check(res, {
        'is status 200': (r) => r.status === 200,
        'code is 0': (r) => JSON.parse(r.body).code === 0,
      });
    }

    // 模拟刷新用户
    {
      sleep(Math.random() * 3)
      let url = `http://${__ENV.MY_HOSTNAME}/app/v2/refresh`
      let params = {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': "bearer " + obtainedJWT
        },
      };
      let payload = JSON.stringify({
        "device": {
          "imei": "string",
          "idfa": "string",
          "android_id": "string",
          "mac": "string",
          "os": "OS_UNKNOWN",
          "oaid": "string",
          "smid": "string",
          "suuid": suuid
        },
        "channel": "performance testing",
        "version_code": "1001"
      });
      let res = http.post(url, payload, params);
      check(res, {
        'is status 200': (r) => r.status === 200,
        'code is 0': (r) => JSON.parse(r.body).code === 0,
      });
    }
  })
  group("rule tests", function () {
    const url =
        `http://${__ENV.MY_HOSTNAME}/rule/v1/calculate/appGlobalCrashConfig-prod?imei=869829049740691&idfa=&android_id=8ea30db01820e2e7&suuid=DoNews5a917964-36ab-4392-9961-2c0324f18087&mac=02:00:00:00:00:00&os=2&user_id=0&oaid=&channel=kuaishou-16&version_code=10036&package_name=com.skin.v10mogul`;
    let res = http.get(url)
    check(res, {
      'is status 200': (r) => r.status === 200,
      'code is 0': (r) => r.json().code === 0,
    });
    sleep(1)
  })

  group("share tests", function () {
    // 填写邀请码
    {
      let url = `http://${__ENV.MY_HOSTNAME}/share/v1/code`;
      let payload = JSON.stringify({
        "invite_code": "0w427ZoMzB"
      });

      let params = {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'bearer ' + JWT()
        },
      };

      let res = http.put(url, payload, params);
      check(res, {
        'is status 200': (r) => r.status === 200,
        'code is 0 or 9': (r) => r.json().code === 0 || r.json().code === 9,
      });
      sleep(1)
    }
  })
}

