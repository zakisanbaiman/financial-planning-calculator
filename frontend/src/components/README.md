# Common Component Library

This directory contains reusable UI components styled with Tailwind CSS.

## Components

### InputField
A flexible input field component with label, error, and helper text support.

**Props:**
- `label?: string` - Label text
- `error?: string` - Error message to display
- `helperText?: string` - Helper text below input
- `fullWidth?: boolean` - Whether input should take full width (default: true)
- All standard HTML input attributes

**Example:**
```tsx
<InputField
  label="月収"
  type="number"
  placeholder="400000"
  error={errors.income}
  helperText="税引き後の手取り額を入力してください"
  required
/>
```

### Button
A versatile button component with multiple variants and sizes.

**Props:**
- `variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'outline'` - Button style (default: 'primary')
- `size?: 'sm' | 'md' | 'lg'` - Button size (default: 'md')
- `fullWidth?: boolean` - Whether button should take full width (default: false)
- `loading?: boolean` - Show loading spinner (default: false)
- All standard HTML button attributes

**Example:**
```tsx
<Button variant="primary" size="md" onClick={handleSubmit} loading={isSubmitting}>
  計算する
</Button>

<Button variant="outline" size="sm" onClick={handleCancel}>
  キャンセル
</Button>
```

### Modal
A modal dialog component with backdrop and keyboard support.

**Props:**
- `isOpen: boolean` - Whether modal is visible
- `onClose: () => void` - Callback when modal should close
- `title?: string` - Modal title
- `size?: 'sm' | 'md' | 'lg' | 'xl'` - Modal size (default: 'md')
- `showCloseButton?: boolean` - Show close button (default: true)
- `children: ReactNode` - Modal content

**Example:**
```tsx
<Modal
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  title="目標を追加"
  size="lg"
>
  <GoalForm onSubmit={handleGoalSubmit} />
</Modal>
```

### LoadingSpinner
A loading spinner component with customizable size and color.

**Props:**
- `size?: 'sm' | 'md' | 'lg' | 'xl'` - Spinner size (default: 'md')
- `color?: 'primary' | 'white' | 'gray'` - Spinner color (default: 'primary')
- `fullScreen?: boolean` - Show as full-screen overlay (default: false)
- `text?: string` - Optional loading text

**Example:**
```tsx
<LoadingSpinner size="lg" text="計算中..." />

<LoadingSpinner fullScreen text="データを読み込んでいます..." />
```

## Usage

Import components from the index file:

```tsx
import { InputField, Button, Modal, LoadingSpinner } from '@/components';
```

Or import individually:

```tsx
import Button from '@/components/Button';
import InputField from '@/components/InputField';
```

## Styling

All components use Tailwind CSS for styling and follow the design system defined in `tailwind.config.js`:

- Primary color: Blue (#3b82f6)
- Success color: Green (#22c55e)
- Warning color: Amber (#f59e0b)
- Error color: Red (#ef4444)

Components are fully responsive and accessible with proper ARIA attributes.
