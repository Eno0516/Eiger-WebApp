import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { Button } from './Button'

describe('Button', () => {
  describe('表示', () => {
    it('children をレンダリングする', () => {
      // Arrange
      // Act
      render(<Button>送信</Button>)
      // Assert
      expect(screen.getByRole('button', { name: '送信' })).toBeInTheDocument()
    })

    it('loading=true のとき "Loading..." と表示され、ボタンが disabled になる', () => {
      // Arrange
      // Act
      render(<Button loading>送信</Button>)
      // Assert
      expect(screen.getByRole('button')).toBeDisabled()
      expect(screen.getByText('Loading...')).toBeInTheDocument()
    })

    it('disabled=true のとき ボタンが disabled になる', () => {
      // Arrange
      // Act
      render(<Button disabled>送信</Button>)
      // Assert
      expect(screen.getByRole('button')).toBeDisabled()
    })
  })

  describe('variant', () => {
    it('variant="primary"（デフォルト）のとき primary スタイルが適用される', () => {
      // Arrange
      // Act
      render(<Button variant="primary">送信</Button>)
      // Assert
      expect(screen.getByRole('button')).toHaveStyle({ background: '#3b82f6' })
    })

    it('variant="secondary" のとき secondary スタイルが適用される', () => {
      // Arrange
      // Act
      render(<Button variant="secondary">キャンセル</Button>)
      // Assert
      expect(screen.getByRole('button')).toHaveStyle({ background: '#e5e7eb' })
    })

    it('variant="danger" のとき danger スタイルが適用される', () => {
      // Arrange
      // Act
      render(<Button variant="danger">削除</Button>)
      // Assert
      expect(screen.getByRole('button')).toHaveStyle({ background: '#ef4444' })
    })
  })

  describe('インタラクション', () => {
    it('クリックで onClick ハンドラが呼ばれる', async () => {
      // Arrange
      const handleClick = vi.fn()
      render(<Button onClick={handleClick}>送信</Button>)
      // Act
      await userEvent.click(screen.getByRole('button'))
      // Assert
      expect(handleClick).toHaveBeenCalledTimes(1)
    })

    it('disabled のとき クリックしても onClick が呼ばれない', async () => {
      // Arrange
      const handleClick = vi.fn()
      render(<Button disabled onClick={handleClick}>送信</Button>)
      // Act
      await userEvent.click(screen.getByRole('button'))
      // Assert
      expect(handleClick).not.toHaveBeenCalled()
    })

    it('loading のとき クリックしても onClick が呼ばれない', async () => {
      // Arrange
      const handleClick = vi.fn()
      render(<Button loading onClick={handleClick}>送信</Button>)
      // Act
      await userEvent.click(screen.getByRole('button'))
      // Assert
      expect(handleClick).not.toHaveBeenCalled()
    })
  })
})
