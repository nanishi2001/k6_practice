import http from 'k6/http';
import { check, sleep } from 'k6';
import { Options } from 'k6/options';

// スパイクテスト: 急激な負荷変動への耐性を確認
export const options: Options = {
  stages: [
    { duration: '30s', target: 5 },    // 通常状態
    { duration: '10s', target: 100 },  // 急激にスパイク
    { duration: '1m', target: 100 },   // スパイク維持
    { duration: '10s', target: 5 },    // 急激に減少
    { duration: '30s', target: 5 },    // 回復確認
    { duration: '10s', target: 0 },    // 終了
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // スパイク時は2秒まで許容
    http_req_failed: ['rate<0.10'],    // エラー率10%未満
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export default function (): void {
  // ランダム遅延エンドポイント（スパイク耐性テスト用）
  const randomDelayRes = http.get(`${BASE_URL}/random-delay`);
  check(randomDelayRes, {
    'random-delay: status is 200': (r) => r.status === 200,
    'random-delay: has body': (r) => r.body !== null && r.body.length > 0,
  });

  // ヘルスチェック
  const healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health: status is 200': (r) => r.status === 200,
  });

  // エラーレートエンドポイント（5%エラー）
  const errorRes = http.get(`${BASE_URL}/error-rate/5`);
  check(errorRes, {
    'error-rate: responded': (r) => r.status === 200 || r.status === 500,
  });

  sleep(0.5);
}
