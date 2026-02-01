package handlers

// テスト用のcredential情報
// 本番環境の値とは異なる、テスト専用の値を使用
const (
	// テスト用パスワード（auth.goのテスト用固定値と一致）
	TestValidPassword   = "password"
	TestInvalidPassword = "wrongpassword"

	// テスト用メールアドレス（初期データのユーザー）
	TestUserEmailAlice   = "alice@example.com"
	TestUserEmailBob     = "bob@example.com"
	TestUserEmailCharlie = "charlie@example.com"
	TestUserEmailUnknown = "unknown@example.com"

	// テスト用ユーザー作成データ
	TestNewUserName  = "Test User"
	TestNewUserEmail = "test@example.com"

	// テスト用更新データ
	TestUpdatedUserName  = "Updated Name"
	TestUpdatedUserEmail = "updated@example.com"

	// 無効なトークン
	TestInvalidToken = "invalid-token"

	// 期待されるトークン有効期限（秒）
	TestExpectedExpiresIn = 900
)
