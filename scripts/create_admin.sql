-- 初期管理者ユーザーを作成
-- パスワード: admin123
-- bcryptハッシュ: cost 12
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYmkjEZvvCa',
    'admin',
    'active'
)
ON DUPLICATE KEY UPDATE username=username;

-- サンプルカテゴリを作成
INSERT INTO categories (name, slug, description)
VALUES
    ('技術', 'tech', '技術関連の記事'),
    ('日記', 'diary', '日常の出来事'),
    ('お知らせ', 'news', 'お知らせ')
ON DUPLICATE KEY UPDATE name=name;

-- サンプルタグを作成
INSERT INTO tags (name, slug)
VALUES
    ('Go', 'go'),
    ('Web開発', 'web-development'),
    ('データベース', 'database'),
    ('セキュリティ', 'security'),
    ('Docker', 'docker')
ON DUPLICATE KEY UPDATE name=name;

