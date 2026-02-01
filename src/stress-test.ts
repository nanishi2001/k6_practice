import http from 'k6/http';
import { check, sleep } from 'k6';
import { Options } from 'k6/options';

// ストレステスト: 段階的に負荷を増加させて限界を特定
export const options: Options = {
  stages: [
    { duration: '2m', target: 10 },   // ウォームアップ
    { duration: '3m', target: 50 },   // 通常負荷
    { duration: '3m', target: 100 },  // 高負荷
    { duration: '3m', target: 150 },  // ストレス負荷
    { duration: '3m', target: 200 },  // 限界負荷
    { duration: '2m', target: 0 },    // クールダウン
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // ストレス時は1秒まで許容
    http_req_failed: ['rate<0.05'],    // エラー率5%未満
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export default function (): void {
  // ユーザー作成
  const createRes = http.post(
    `${BASE_URL}/users`,
    JSON.stringify({ name: `user_${Date.now()}`, email: `test_${Date.now()}@example.com` }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(createRes, {
    'create: status is 201': (r) => r.status === 201,
    'create: response time < 1000ms': (r) => r.timings.duration < 1000,
  });

  // ユーザー一覧取得
  const listRes = http.get(`${BASE_URL}/users`);
  check(listRes, {
    'list: status is 200': (r) => r.status === 200,
  });

  // 遅延エンドポイント（負荷テスト用）
  const delayRes = http.get(`${BASE_URL}/delay/50`);
  check(delayRes, {
    'delay: status is 200': (r) => r.status === 200,
  });

  sleep(Math.random() * 2 + 1); // 1-3秒のランダムスリープ
}
