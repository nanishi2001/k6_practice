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

const BASE_URL = 'https://test.k6.io';

export default function (): void {
  const response = http.get(BASE_URL);

  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
