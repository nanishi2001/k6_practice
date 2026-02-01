import http from 'k6/http';
import { check, sleep } from 'k6';
import { Options } from 'k6/options';

// 負荷テスト: 一定のユーザー数で通常負荷をシミュレート
export const options: Options = {
  stages: [
    { duration: '1m', target: 10 },  // 1分で10VUsまで増加
    { duration: '3m', target: 10 },  // 3分間10VUsを維持
    { duration: '1m', target: 0 },   // 1分でクールダウン
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95%のリクエストが500ms以下
    http_req_failed: ['rate<0.01'],    // エラー率1%未満
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export default function (): void {
  // ヘルスチェック
  const healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health: status is 200': (r) => r.status === 200,
    'health: response time < 500ms': (r) => r.timings.duration < 500,
  });

  // ユーザー一覧取得
  const usersRes = http.get(`${BASE_URL}/users`);
  check(usersRes, {
    'users: status is 200': (r) => r.status === 200,
    'users: response time < 500ms': (r) => r.timings.duration < 500,
  });

  // 個別ユーザー取得
  const userRes = http.get(`${BASE_URL}/users/1`);
  check(userRes, {
    'user: status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
