import {uuidv4} from "https://jslib.k6.io/k6-utils/1.0.0/index.js";
import crypto from "k6/crypto";
import encoding from "k6/encoding";

const algToHash = {
  HS256: "sha256",
  HS384: "sha384",
  HS512: "sha512"
};

function sign(data, hashAlg, secret) {
  let hasher = crypto.createHMAC(hashAlg, secret);
  hasher.update(data);

  // Some manual base64 rawurl encoding as `Hasher.digest(encodingType)`
  // doesn't support that encoding type yet.
  return hasher.digest("base64").replace(/\//g, "_").replace(/\+/g,
      "-").replace(/=/g, "");
}

export function encode(payload, secret, algorithm) {
  algorithm = algorithm || "HS256";
  let header = encoding.b64encode(JSON.stringify({typ: "JWT", alg: algorithm}),
      "rawurl");
  payload = encoding.b64encode(JSON.stringify(payload), "rawurl");
  let sig = sign(header + "." + payload, algToHash[algorithm], secret);
  return [header, payload, sig].join(".");
}

export function decode(token, secret, algorithm) {
  let parts = token.split('.');
  let header = JSON.parse(encoding.b64decode(parts[0], "rawurl"));
  let payload = JSON.parse(encoding.b64decode(parts[1], "rawurl"));
  algorithm = algorithm || algToHash[header.alg];
  if (sign(parts[0] + "." + parts[1], algorithm, secret) != parts[2]) {
    throw Error("JWT signature verification failed");
  }
  return payload;
}

function getRandomInt(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min)) + min; //不含最大值，含最小值
}

export function JWT(option = {}) {
  let now = Math.floor(Date.now() / 1000)
  let nextYear = Math.floor(new Date(new Date().setFullYear(new Date().getFullYear() + 1)) / 1000)
  let payload = {
    PackageName: option.PackageName || "",
    UserId: option.UserId || getRandomInt(1, 100000),
    Suuid: option.Suuid || uuidv4(),
    Channel: option.Channel || "performance testing",
    VersionCode: option.VersionCode || "1000",
    Wechat: option.Wechat || "",
    Mobile: option.Mobile || "",
    ThirdPartyId: option.ThirdPartyId || "",
    aud: option.aud || "",
    exp: option.exp || nextYear,
    jti: option.id || getRandomInt(1, 100000).toString(),
    iat: option.iat || now,
    iss: option.iss || "performance testing",
    nbf: option.nbf || now,
    sub: option.sub || "",
  }
  let secret = "zxcvb0997zSDvHSD"
  let algorithm = "HS256"
  let token = encode(payload, secret, algorithm)
  return token
}
