import http from 'k6/http';
import {check, sleep} from 'k6';
import {URL} from 'https://jslib.k6.io/url/1.0.0/index.js';

export let options = {
  vus: 1000,
  duration: '15s',
  thresholds: {
    http_req_failed: ['rate<0.01'],   // http errors should be less than 1%
    http_req_duration: ['p(90)<100'], // 95% of requests should be below 200ms
  },
};

export default function () {
  const url = new URL('https://monetization.tagtic.cn/rule/v1/calculate/appGlobalCrashConfig-prod?imei=869829049740691&idfa=&android_id=8ea30db01820e2e7&suuid=DoNews5a917964-36ab-4392-9961-2c0324f18087&mac=02:00:00:00:00:00&os=2&user_id=0&oaid=&channel=kuaishou-16&version_code=10036&package_name=com.skin.v10mogul');
  let res = http.get(url.toString())
  check(res, {
    'is status 200': (r) => r.status === 200,
    'code is 0': (r) => r.json().code === 0,
  });
  sleep(1)
}
