// Common UI Components
export { default as InputField } from './InputField';
export type { InputFieldProps } from './InputField';

export { default as Button } from './Button';
export type { ButtonProps, ButtonVariant, ButtonSize } from './Button';

export { default as Modal } from './Modal';
export type { ModalProps } from './Modal';

export { default as LoadingSpinner } from './LoadingSpinner';
export type { LoadingSpinnerProps, SpinnerSize, SpinnerColor } from './LoadingSpinner';

export { default as Navigation } from './Navigation';

export { default as ThemeToggle } from './ThemeToggle';

export { default as Tutorial } from './Tutorial';

// Form Components
export { default as FinancialInputForm } from './FinancialInputForm';
export type { FinancialInputFormProps } from './FinancialInputForm';

export { default as InvestmentSettingsForm } from './InvestmentSettingsForm';
export type { InvestmentSettingsFormProps, InvestmentSettings } from './InvestmentSettingsForm';

export { default as GoalForm } from './GoalForm';
export type { GoalFormProps } from './GoalForm';

export { default as GoalProgressTracker } from './GoalProgressTracker';
export type { GoalProgressTrackerProps } from './GoalProgressTracker';

export { default as GoalsSummaryChart } from './GoalsSummaryChart';
export type { GoalsSummaryChartProps } from './GoalsSummaryChart';

export { default as GoalRecommendations } from './GoalRecommendations';
export type { GoalRecommendationsProps } from './GoalRecommendations';

export { default as CurrencyInput } from './CurrencyInput';
export type { CurrencyInputProps } from './CurrencyInput';

export { default as CurrencyInputWithPresets } from './CurrencyInputWithPresets';
export type { CurrencyInputWithPresetsProps, PresetValue } from './CurrencyInputWithPresets';

// Calculation Components
export { default as AssetProjectionChart } from './AssetProjectionChart';
export { default as AssetProjectionCalculator } from './AssetProjectionCalculator';
export { default as RetirementCalculator } from './RetirementCalculator';
export { default as EmergencyFundCalculator } from './EmergencyFundCalculator';

// Error Handling Components
export { ErrorBoundary, APIErrorDisplay } from './ErrorBoundary';

// Connection Status Components
export { ConnectionStatus, InlineConnectionStatus } from './ConnectionStatus';
