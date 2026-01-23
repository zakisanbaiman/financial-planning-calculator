import React, { forwardRef, InputHTMLAttributes } from 'react';

export interface InputFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  fullWidth?: boolean;
}

const InputField = forwardRef<HTMLInputElement, InputFieldProps>(
  ({ label, error, helperText, fullWidth = true, className = '', ...props }, ref) => {
    const inputClasses = `
      px-3 py-3 border rounded-lg text-base
      focus:outline-none focus:ring-2 focus:border-transparent
      disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed
      dark:bg-gray-800 dark:text-white
      ${error ? 'border-error-500 focus:ring-error-500' : 'border-gray-300 dark:border-gray-600 focus:ring-primary-500'}
      ${fullWidth ? 'w-full' : ''}
      ${className}
    `.trim().replace(/\s+/g, ' ');

    return (
      <div className={fullWidth ? 'w-full' : ''}>
        {label && (
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {label}
            {props.required && <span className="text-error-500 ml-1">*</span>}
          </label>
        )}
        <input
          ref={ref}
          className={inputClasses}
          style={{ fontSize: '16px' }} // iOS Safariのズームを防ぐ
          {...props}
        />
        {error && (
          <p className="mt-2 text-sm text-error-500">{error}</p>
        )}
        {helperText && !error && (
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">{helperText}</p>
        )}
      </div>
    );
  }
);

InputField.displayName = 'InputField';

export default InputField;
