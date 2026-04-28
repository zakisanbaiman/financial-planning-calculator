const nextConfig = require('eslint-config-next/core-web-vitals');

module.exports = [
  ...nextConfig,
  {
    rules: {
      // Downgrade React Compiler rules introduced in Next.js 16 to warnings
      // for gradual adoption on existing code
      'react-hooks/set-state-in-effect': 'warn',
      'react-hooks/set-state-in-render': 'warn',
    },
  },
];
