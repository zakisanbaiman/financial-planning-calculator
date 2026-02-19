'use client';

import React, { useState, useEffect, forwardRef } from 'react';

export interface CurrencyInputProps {
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
}

/**
 * 金額入力コンポーネント
 * 入力中は数値のみ、フォーカスが外れると3桁区切りで表示
 */
const CurrencyInput = forwardRef<HTMLInputElement, CurrencyInputProps>(
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
    },
    ref
  ) => {
    const [isFocused, setIsFocused] = useState(false);
    const [inputValue, setInputValue] = useState('');

    // 値が変更されたときに表示を更新
    useEffect(() => {
      if (!isFocused) {
        setInputValue(value ? value.toLocaleString() : '');
      }
    }, [value, isFocused]);

    const handleFocus = () => {
      setIsFocused(true);
      // フォーカス時もカンマ付きのまま表示
    };

    const handleBlur = () => {
      setIsFocused(false);
      // フォーカスが外れても3桁区切りで表示（変更なし）
      setInputValue(value ? value.toLocaleString() : '');
      onBlur?.();
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const input = e.target;
      const cursorPos = input.selectionStart ?? 0;
      const oldFormattedValue = input.value;

      // カーソル前のカンマ数・数字数を取得
      const oldCommasBeforeCursor = (oldFormattedValue.slice(0, cursorPos).match(/,/g) || []).length;
      const digitsBeforeCursor = cursorPos - oldCommasBeforeCursor;

      const rawValue = oldFormattedValue.replace(/,/g, ''); // カンマを除去

      // 空文字の場合は0
      if (rawValue === '') {
        setInputValue('');
        onChange(0);
        return;
      }

      // 数値のみ許可
      if (!/^\d*$/.test(rawValue)) {
        return;
      }

      const numValue = parseInt(rawValue, 10);

      // min/max チェック
      if (min !== undefined && numValue < min) {
        return;
      }
      if (max !== undefined && numValue > max) {
        return;
      }

      // 入力中も3桁区切りでフォーマット
      const formatted = numValue.toLocaleString();
      setInputValue(formatted);
      onChange(numValue);

      // カーソル位置をカンマ挿入後も正しい位置に調整
      let digitCount = 0;
      let newCursorPos = formatted.length;

      if (digitsBeforeCursor === 0) {
        newCursorPos = 0;
      } else {
        for (let i = 0; i < formatted.length; i++) {
          if (formatted[i] !== ',') {
            digitCount++;
          }
          if (digitCount === digitsBeforeCursor) {
            newCursorPos = i + 1;
            break;
          }
        }
      }

      requestAnimationFrame(() => {
        if (input.isConnected) {
          input.setSelectionRange(newCursorPos, newCursorPos);
        }
      });
    };

    const inputClasses = `
      px-3 py-2 border rounded-lg w-full
      focus:outline-none focus:ring-2 focus:border-transparent
      disabled:bg-gray-100 disabled:cursor-not-allowed
      ${error ? 'border-error-500 focus:ring-error-500' : 'border-gray-300 focus:ring-primary-500'}
    `.trim().replace(/\s+/g, ' ');

    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {label}
            {required && <span className="text-error-500 ml-1">*</span>}
          </label>
        )}
        <div className="relative">
          <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">
            ¥
          </span>
          <input
            ref={ref}
            type="text"
            inputMode="numeric"
            className={`${inputClasses} pl-7`}
            value={inputValue}
            onChange={handleChange}
            onFocus={handleFocus}
            onBlur={handleBlur}
            placeholder={placeholder}
            disabled={disabled}
          />
        </div>
        {error && <p className="mt-1 text-sm text-error-500">{error}</p>}
        {helperText && !error && <p className="mt-1 text-sm text-gray-500">{helperText}</p>}
      </div>
    );
  }
);

CurrencyInput.displayName = 'CurrencyInput';

export default CurrencyInput;
