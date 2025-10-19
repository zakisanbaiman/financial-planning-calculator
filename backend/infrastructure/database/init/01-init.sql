-- PostgreSQL初期化スクリプト
-- Docker起動時に自動実行される

-- UUIDエクステンションを有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 日本語ロケール設定（必要に応じて）
-- CREATE COLLATION IF NOT EXISTS japanese (locale = 'ja_JP.UTF-8');

-- データベースの基本設定
ALTER DATABASE financial_planning SET timezone TO 'Asia/Tokyo';

-- 接続確認用のテーブル（マイグレーションで削除される）
CREATE TABLE IF NOT EXISTS _docker_init_check (
    id SERIAL PRIMARY KEY,
    initialized_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO _docker_init_check (initialized_at) VALUES (CURRENT_TIMESTAMP);