'use client';

import React, { forwardRef } from 'react';

export interface PresetValue {
  label: string;
  value: number;
}

export interface CurrencyInputWithPresetsProps {
  label?: string;
  value: number;
  onChange: (value: number) => void;
  onBlur?: () => void;
  error?: string;
  helperText?: string;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  min?: number;
  max?: number;
  presets?: PresetValue[];
  unit?: string;
}

/**
 * プリセットボタン付き通貨入力コンポーネント
 * モバイルでの大きな数値入力を簡単にするためのプリセットボタンを提供
 */
const CurrencyInputWithPresets = forwardRef<HTMLInputElement, CurrencyInputWithPresetsProps>(
  (
    {
      label,
      value,
      onChange,
      onBlur,
      error,
      helperText,
      placeholder = '0',
      required = false,
      disabled = false,
      min,
      max,
      presets = [],
      unit = '円',
    },
    ref
  ) => {
    const handlePresetClick = (presetValue: number) => {
      if (disabled) return;
      onChange(presetValue);
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const rawValue = e.target.value.replace(/,/g, '');
      
      if (rawValue === '') {
        onChange(0);
        return;
      }

      if (!/^\d*$/.test(rawValue)) {
        return;
      }

      const numValue = parseInt(rawValue, 10);
      
      if (min !== undefined && numValue < min) {
        return;
      }
      if (max !== undefined && numValue > max) {
        return;
      }

      onChange(numValue);
    };

    const displayValue = value ? value.toLocaleString() : '';

    const inputClasses = `
      px-3 py-3 border rounded-lg w-full text-base
      focus:outline-none focus:ring-2 focus:border-transparent
      disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed
      dark:bg-gray-800 dark:text-white
      ${error ? 'border-error-500 focus:ring-error-500' : 'border-gray-300 dark:border-gray-600 focus:ring-primary-500'}
    `.trim().replace(/\s+/g, ' ');

    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {label}
            {required && <span className="text-error-500 ml-1">*</span>}
          </label>
        )}
        
        {/* プリセットボタン */}
        {presets.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-3">
            {presets.map((preset) => (
              <button
                key={preset.value}
                type="button"
                onClick={() => handlePresetClick(preset.value)}
                disabled={disabled}
                className={`
                  px-3 py-2 text-sm font-medium rounded-lg transition-colors
                  ${
                    value === preset.value
                      ? 'bg-primary-500 text-white'
                      : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                  }
                  disabled:opacity-50 disabled:cursor-not-allowed
                  active:scale-95 transform
                  min-w-[44px] min-h-[44px] flex items-center justify-center
                `}
              >
                {preset.label}
              </button>
            ))}
          </div>
        )}

        {/* 入力フィールド */}
        <div className="relative">
          <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-400 text-base">
            ¥
          </span>
          <input
            ref={ref}
            type="text"
            inputMode="numeric"
            className={`${inputClasses} pl-8 pr-12`}
            value={displayValue}
            onChange={handleInputChange}
            onBlur={onBlur}
            placeholder={placeholder}
            disabled={disabled}
            style={{ fontSize: '16px' }} // iOS Safariのズームを防ぐ
          />
          <span className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-400 text-sm">
            {unit}
          </span>
        </div>

        {error && <p className="mt-2 text-sm text-error-500">{error}</p>}
        {helperText && !error && <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">{helperText}</p>}
      </div>
    );
  }
);

CurrencyInputWithPresets.displayName = 'CurrencyInputWithPresets';

export default CurrencyInputWithPresets;
