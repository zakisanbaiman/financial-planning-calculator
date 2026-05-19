const nextConfig = require('eslint-config-next/core-web-vitals');

module.exports = [
  ...nextConfig,
  {
    rules: {
      // React Compiler rules introduced in Next.js 16 – disabled for existing
      // codebase that predates React Compiler adoption
      'react-hooks/set-state-in-effect': 'off',
      'react-hooks/set-state-in-render': 'off',
      'react-hooks/incompatible-library': 'off',
    },
  },
];
