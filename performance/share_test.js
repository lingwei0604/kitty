import {check, sleep} from 'k6';
import http from 'k6/http';
import {JWT} from './jwt.js';

export let options = {
  vus: 1000,
  duration: '15s',
  thresholds: {
    http_req_failed: ['rate<0.01'],   // http errors should be less than 1%
    http_req_duration: ['p(90)<500'], // 95% of requests should be below 200ms
  },
};

export default function () {
  // 填写邀请码
  {
    let url = 'https://monetization.dev.tagtic.cn/share/v1/code';
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
}
