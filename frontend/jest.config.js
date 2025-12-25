const nextJest = require('next/jest');

const createJestConfig = nextJest({
  // next.config.js と .env ファイルを読み込むための Next.js アプリのパスを指定
  dir: './',
});

// Jest に渡すカスタム設定
const customJestConfig = {
  // テストファイルのパターン
  testMatch: [
    '**/__tests__/**/*.test.[jt]s?(x)',
    '**/__tests__/**/*.spec.[jt]s?(x)',
  ],
  
  // セットアップファイル
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  
  // テスト環境
  testEnvironment: 'jest-environment-jsdom',
  
  // モジュールパスエイリアス（tsconfig.json の paths と一致させる）
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  
  // カバレッジ設定
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.d.ts',
    '!src/**/index.ts',
    '!src/**/*.stories.{js,jsx,ts,tsx}',
    '!src/app/layout.tsx',
    '!src/app/page.tsx',
  ],
  
  // カバレッジ閾値
  coverageThreshold: {
    global: {
      branches: 50,
      functions: 50,
      lines: 50,
      statements: 50,
    },
  },
  
  // テスト対象外
  testPathIgnorePatterns: [
    '<rootDir>/node_modules/',
    '<rootDir>/.next/',
  ],
  
  // TypeScript 変換
  transform: {
    '^.+\\.(ts|tsx)$': ['ts-jest', {
      tsconfig: 'tsconfig.json',
    }],
  },
  
  // モジュール拡張子
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json'],
};

// createJestConfig は非同期で next/jest が Next.js 設定を読み込めるようにエクスポートされる
module.exports = createJestConfig(customJestConfig);
