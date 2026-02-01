---
name: k6-new
description: 新しいk6テストスクリプトのテンプレートを生成
user_invocable: true
---

# /k6-new

新しいk6テストスクリプトを作成します。

## 使用方法

```
/k6-new <name> [type]
```

- `name` - テスト名（例: `api-health`）
- `type` - テストタイプ（load/stress/spike/soak）

## テンプレート

```typescript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Options } from 'k6/options';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options: Options = {
  stages: [
    { duration: '1m', target: 10 },
    { duration: '3m', target: 10 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function (): void {
  const res = http.get(`${BASE_URL}/endpoint`);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
```
