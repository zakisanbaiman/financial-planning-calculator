-- 001_sample_data.sql
-- 開発・テスト用のサンプルデータ

-- サンプルユーザーの作成
INSERT INTO users (id, email) VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'user1@example.com'),
    ('550e8400-e29b-41d4-a716-446655440002', 'user2@example.com')
ON CONFLICT (email) DO NOTHING;

-- ユーザー1の財務データ
INSERT INTO financial_data (
    id, user_id, monthly_income, investment_return, inflation_rate
) VALUES (
    '550e8400-e29b-41d4-a716-446655440011',
    '550e8400-e29b-41d4-a716-446655440001',
    400000.00,
    5.0,
    2.0
) ON CONFLICT (user_id) DO NOTHING;

-- ユーザー1の支出項目
INSERT INTO expense_items (financial_data_id, category, amount, description) VALUES 
    ('550e8400-e29b-41d4-a716-446655440011', '住居費', 120000.00, '家賃・管理費'),
    ('550e8400-e29b-41d4-a716-446655440011', '食費', 60000.00, '食材・外食費'),
    ('550e8400-e29b-41d4-a716-446655440011', '交通費', 20000.00, '通勤・交通費'),
    ('550e8400-e29b-41d4-a716-446655440011', '光熱費', 15000.00, '電気・ガス・水道'),
    ('550e8400-e29b-41d4-a716-446655440011', '通信費', 12000.00, '携帯・インターネット'),
    ('550e8400-e29b-41d4-a716-446655440011', '保険料', 25000.00, '生命保険・医療保険'),
    ('550e8400-e29b-41d4-a716-446655440011', 'その他', 48000.00, '娯楽・雑費')
ON CONFLICT DO NOTHING;

-- ユーザー1の貯蓄項目
INSERT INTO savings_items (financial_data_id, type, amount, description) VALUES 
    ('550e8400-e29b-41d4-a716-446655440011', 'deposit', 1000000.00, '普通預金'),
    ('550e8400-e29b-41d4-a716-446655440011', 'deposit', 500000.00, '定期預金'),
    ('550e8400-e29b-41d4-a716-446655440011', 'investment', 800000.00, '投資信託'),
    ('550e8400-e29b-41d4-a716-446655440011', 'investment', 300000.00, '株式投資')
ON CONFLICT DO NOTHING;

-- ユーザー1の退職データ
INSERT INTO retirement_data (
    id, user_id, current_age, retirement_age, life_expectancy, 
    monthly_retirement_expenses, pension_amount
) VALUES (
    '550e8400-e29b-41d4-a716-446655440021',
    '550e8400-e29b-41d4-a716-446655440001',
    35,
    65,
    85,
    250000.00,
    150000.00
) ON CONFLICT (user_id) DO NOTHING;

-- ユーザー1の目標
INSERT INTO goals (user_id, type, title, target_amount, target_date, current_amount, monthly_contribution) VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'emergency', '緊急資金', 1500000.00, '2025-12-31', 500000.00, 50000.00),
    ('550e8400-e29b-41d4-a716-446655440001', 'savings', 'マイホーム頭金', 5000000.00, '2028-03-31', 800000.00, 100000.00),
    ('550e8400-e29b-41d4-a716-446655440001', 'retirement', '老後資金', 30000000.00, '2054-12-31', 1100000.00, 80000.00),
    ('550e8400-e29b-41d4-a716-446655440001', 'custom', '子供の教育資金', 3000000.00, '2035-03-31', 200000.00, 60000.00)
ON CONFLICT DO NOTHING;

-- ユーザー2の財務データ
INSERT INTO financial_data (
    id, user_id, monthly_income, investment_return, inflation_rate
) VALUES (
    '550e8400-e29b-41d4-a716-446655440012',
    '550e8400-e29b-41d4-a716-446655440002',
    600000.00,
    6.0,
    2.5
) ON CONFLICT (user_id) DO NOTHING;

-- ユーザー2の支出項目
INSERT INTO expense_items (financial_data_id, category, amount, description) VALUES 
    ('550e8400-e29b-41d4-a716-446655440012', '住居費', 180000.00, '住宅ローン'),
    ('550e8400-e29b-41d4-a716-446655440012', '食費', 80000.00, '食材・外食費'),
    ('550e8400-e29b-41d4-a716-446655440012', '交通費', 30000.00, '車両費・ガソリン'),
    ('550e8400-e29b-41d4-a716-446655440012', '光熱費', 20000.00, '電気・ガス・水道'),
    ('550e8400-e29b-41d4-a716-446655440012', '通信費', 15000.00, '携帯・インターネット'),
    ('550e8400-e29b-41d4-a716-446655440012', '保険料', 35000.00, '生命保険・医療保険・車両保険'),
    ('550e8400-e29b-41d4-a716-446655440012', '教育費', 50000.00, '子供の習い事・塾'),
    ('550e8400-e29b-41d4-a716-446655440012', 'その他', 90000.00, '娯楽・雑費')
ON CONFLICT DO NOTHING;

-- ユーザー2の貯蓄項目
INSERT INTO savings_items (financial_data_id, type, amount, description) VALUES 
    ('550e8400-e29b-41d4-a716-446655440012', 'deposit', 2000000.00, '普通預金'),
    ('550e8400-e29b-41d4-a716-446655440012', 'investment', 1500000.00, 'つみたてNISA'),
    ('550e8400-e29b-41d4-a716-446655440012', 'investment', 800000.00, 'iDeCo'),
    ('550e8400-e29b-41d4-a716-446655440012', 'other', 500000.00, '学資保険')
ON CONFLICT DO NOTHING;

-- ユーザー2の退職データ
INSERT INTO retirement_data (
    id, user_id, current_age, retirement_age, life_expectancy, 
    monthly_retirement_expenses, pension_amount
) VALUES (
    '550e8400-e29b-41d4-a716-446655440022',
    '550e8400-e29b-41d4-a716-446655440002',
    42,
    60,
    85,
    300000.00,
    180000.00
) ON CONFLICT (user_id) DO NOTHING;

-- ユーザー2の目標
INSERT INTO goals (user_id, type, title, target_amount, target_date, current_amount, monthly_contribution) VALUES 
    ('550e8400-e29b-41d4-a716-446655440002', 'emergency', '緊急資金', 2400000.00, '2025-06-30', 1000000.00, 70000.00),
    ('550e8400-e29b-41d4-a716-446655440002', 'custom', '車の買い替え', 3500000.00, '2027-12-31', 500000.00, 80000.00),
    ('550e8400-e29b-41d4-a716-446655440002', 'retirement', '早期退職資金', 50000000.00, '2042-12-31', 2300000.00, 150000.00)
ON CONFLICT DO NOTHING;