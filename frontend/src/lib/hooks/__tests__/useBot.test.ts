import { renderHook, act } from '@testing-library/react';
import { useBot } from '../useBot';

// api-client をモック
jest.mock('../../api-client', () => ({
  streamBotAnswer: jest.fn(),
}));

import { streamBotAnswer } from '../../api-client';

const mockStreamBotAnswer = streamBotAnswer as jest.MockedFunction<typeof streamBotAnswer>;

describe('useBot', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('初期状態', () => {
    it('初期値が正しく設定されている', () => {
      const { result } = renderHook(() => useBot());

      expect(result.current.answer).toBe('');
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe('正常系', () => {
    it('ストリーム受信・トークン連結・done後にisLoading=false', async () => {
      // streamBotAnswer がトークンを順に受信して onDone を呼ぶ
      mockStreamBotAnswer.mockImplementation(async (_question, onToken, onDone, _onError) => {
        onToken('こんにちは');
        onToken('！');
        onDone();
      });

      const { result } = renderHook(() => useBot());

      await act(async () => {
        await result.current.ask('挨拶して');
      });

      expect(result.current.answer).toBe('こんにちは！');
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('ask 呼び出し中は isLoading=true になる', async () => {
      let resolveStream: () => void;
      mockStreamBotAnswer.mockImplementation(
        (_question, _onToken, onDone, _onError) =>
          new Promise<void>((resolve) => {
            resolveStream = () => {
              onDone();
              resolve();
            };
          }),
      );

      const { result } = renderHook(() => useBot());

      let askPromise: Promise<void>;
      act(() => {
        askPromise = result.current.ask('質問');
      });

      // まだ解決していないので isLoading=true のはず
      expect(result.current.isLoading).toBe(true);

      await act(async () => {
        resolveStream!();
        await askPromise!;
      });

      expect(result.current.isLoading).toBe(false);
    });
  });

  describe('異常系', () => {
    it('errorイベント受信時にerrorが設定されisLoading=falseになる', async () => {
      mockStreamBotAnswer.mockImplementation(async (_question, _onToken, _onDone, onError) => {
        onError('回答の生成中にエラーが発生しました');
      });

      const { result } = renderHook(() => useBot());

      await act(async () => {
        await result.current.ask('質問');
      });

      expect(result.current.error).toBe('回答の生成中にエラーが発生しました');
      expect(result.current.isLoading).toBe(false);
      expect(result.current.answer).toBe('');
    });

    it('streamBotAnswer が例外をスローした場合にerrorが設定される', async () => {
      mockStreamBotAnswer.mockRejectedValue(new Error('ネットワークエラー'));

      const { result } = renderHook(() => useBot());

      await act(async () => {
        await result.current.ask('質問');
      });

      expect(result.current.error).toBe('ネットワークエラー');
      expect(result.current.isLoading).toBe(false);
    });
  });

  describe('エッジケース', () => {
    it('空質問での送信防止: streamBotAnswer が呼ばれずerrorが設定される', async () => {
      const { result } = renderHook(() => useBot());

      await act(async () => {
        await result.current.ask('');
      });

      expect(mockStreamBotAnswer).not.toHaveBeenCalled();
      expect(result.current.error).toBe('質問を入力してください');
      expect(result.current.isLoading).toBe(false);
    });

    it('空白のみの質問での送信防止', async () => {
      const { result } = renderHook(() => useBot());

      await act(async () => {
        await result.current.ask('   ');
      });

      expect(mockStreamBotAnswer).not.toHaveBeenCalled();
      expect(result.current.error).toBe('質問を入力してください');
    });

    it('reset() を呼ぶと状態が初期値に戻る', async () => {
      mockStreamBotAnswer.mockImplementation(async (_question, _onToken, _onDone, onError) => {
        onError('エラー');
      });

      const { result } = renderHook(() => useBot());

      await act(async () => {
        await result.current.ask('質問');
      });

      expect(result.current.error).not.toBeNull();

      act(() => {
        result.current.reset();
      });

      expect(result.current.answer).toBe('');
      expect(result.current.error).toBeNull();
      expect(result.current.isLoading).toBe(false);
    });
  });
});
