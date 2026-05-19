import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { Input } from './Input'

describe('Input', () => {
  describe('表示', () => {
    it('input 要素をレンダリングする', () => {
      // Arrange
      // Act
      render(<Input />)
      // Assert
      expect(screen.getByRole('textbox')).toBeInTheDocument()
    })

    it('label を指定したとき label 要素が表示される', () => {
      // Arrange
      // Act
      render(<Input label="メールアドレス" id="email" />)
      // Assert
      expect(screen.getByLabelText('メールアドレス')).toBeInTheDocument()
    })

    it('label を指定しないとき label 要素は表示されない', () => {
      // Arrange
      // Act
      render(<Input />)
      // Assert
      expect(document.querySelector('label')).toBeNull()
    })

    it('error を指定したとき エラーメッセージが表示される', () => {
      // Arrange
      // Act
      render(<Input error="必須入力です" />)
      // Assert
      expect(screen.getByText('必須入力です')).toBeInTheDocument()
    })

    it('error を指定したとき input のボーダーが赤くなる', () => {
      // Arrange
      // Act
      render(<Input error="エラー" />)
      // Assert
      expect(screen.getByRole('textbox')).toHaveStyle({ border: '1px solid #ef4444' })
    })

    it('error を指定しないとき エラーメッセージは表示されない', () => {
      // Arrange
      // Act
      render(<Input />)
      // Assert
      expect(document.querySelector('span')).toBeNull()
    })
  })

  describe('インタラクション', () => {
    it('テキストを入力すると onChange が呼ばれる', async () => {
      // Arrange
      const handleChange = vi.fn()
      render(<Input onChange={handleChange} />)
      // Act
      await userEvent.type(screen.getByRole('textbox'), 'test')
      // Assert
      expect(handleChange).toHaveBeenCalled()
    })
  })
})
