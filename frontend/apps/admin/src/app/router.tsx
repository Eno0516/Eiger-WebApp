import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AdminLayout } from '@app/layouts/AdminLayout'
import { CirclesPage } from '../pages/CirclesPage'
import { CircleDetailPage } from '../pages/CircleDetailPage'
import { EventsPage } from '../pages/EventsPage'
import { EventDetailPage } from '../pages/EventDetailPage'

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        {/* 管理画面トップ（サークル一覧ナビ付き） */}
        <Route element={<AdminLayout />}>
          <Route path="/admin" element={<div>SuperAdmin ダッシュボード（未実装）</div>} />
          <Route path="/admin/circles" element={<CirclesPage />} />
          <Route path="/admin/circles/new" element={<div>サークル作成（未実装）</div>} />
        </Route>

        {/* サークル詳細：別タブで開く独立レイアウト */}
        <Route path="/admin/circles/:id" element={<CircleDetailPage />}>
          <Route index element={<Navigate to="members" replace />} />
          <Route path="members" element={<div>メンバー管理（未実装）</div>} />
          <Route path="events" element={<EventsPage />} />
          <Route path="events/new" element={<div>イベント作成（未実装）</div>} />
          <Route path="events/:eventId" element={<EventDetailPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
