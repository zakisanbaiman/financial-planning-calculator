import Link from 'next/link';

const Footer = () => {
  return (
    <footer className="border-t border-ink-200 dark:border-ink-800 bg-ink-50 dark:bg-ink-950 transition-colors mt-auto">
      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <p className="text-sm font-body text-ink-400 dark:text-ink-500">
            © 2026 FinPlan
          </p>
          <nav className="flex items-center gap-6">
            <Link
              href="/terms"
              className="text-sm font-body text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200 transition-colors"
            >
              利用規約
            </Link>
            <Link
              href="/privacy"
              className="text-sm font-body text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200 transition-colors"
            >
              プライバシーポリシー
            </Link>
            <a
              href="https://github.com/zakisanbaiman/financial-planning-calculator"
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm font-body text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200 transition-colors"
            >
              GitHub
            </a>
          </nav>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
