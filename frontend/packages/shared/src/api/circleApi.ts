import { httpClient } from './httpClient'
import type { Circle, CircleMember, FormResponse } from '../types/circle.types'

export const circleApi = {
  // サークル一覧取得
  getAll: () =>
    httpClient.get<Circle[]>('/circles').then((r) => r.data),

  // サークル新規作成（system_admin のみ）
  create: (payload: Omit<Circle, 'id' | 'createdAt' | 'updatedAt'>) =>
    httpClient.post<Circle>('/circles', payload).then((r) => r.data),

  // 参加申請（visitor レコード作成）
  apply: (circleId: string) =>
    httpClient.post<{ message: string }>(`/circles/${circleId}/apply`).then((r) => r.data),

  // メンバー一覧取得（admin のみ）
  getMembers: (circleId: string) =>
    httpClient.get<CircleMember[]>(`/circles/${circleId}/members`).then((r) => r.data),

  // 対象ユーザーの Form 回答を Spreadsheet から取得（admin のみ）
  getFormResponse: (circleId: string, userId: string) =>
    httpClient
      .get<FormResponse>(`/circles/${circleId}/members/${userId}/form-response`)
      .then((r) => r.data),

  // ロール変更（admin のみ）
  updateMemberRole: (circleId: string, userId: string, role: string) =>
    httpClient
      .patch<{ message: string }>(`/circles/${circleId}/members/${userId}/role`, { role })
      .then((r) => r.data),
}
