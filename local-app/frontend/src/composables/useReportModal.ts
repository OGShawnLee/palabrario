import { ref } from 'vue'
import { api } from '@/api/client'
import { errorMessage } from '@/utils/errors'
import type { CreateReportRequest } from '@palabrario/ui'

export function useReportModal() {
  const isReportOpen = ref(false)
  const reportTargetType = ref<'WORD' | 'SENSE' | 'EXAMPLE'>('WORD')
  const reportTargetId = ref(0)
  const reportTargetLabel = ref('')
  const isReportSubmitting = ref(false)
  const reportFeedbackMsg = ref<string | null>(null)

  function showFeedback(msg: string, duration = 4500) {
    reportFeedbackMsg.value = msg
    setTimeout(() => {
      if (reportFeedbackMsg.value === msg) {
        reportFeedbackMsg.value = null
      }
    }, duration)
  }

  function openReport(
    targetType: 'WORD' | 'SENSE' | 'EXAMPLE',
    targetId: number,
    targetLabel: string,
  ) {
    reportTargetType.value = targetType
    reportTargetId.value = targetId
    reportTargetLabel.value = targetLabel
    isReportOpen.value = true
  }

  async function handleReportSubmit(payload: CreateReportRequest): Promise<boolean> {
    isReportSubmitting.value = true
    try {
      await api.CreateReport(payload)
      isReportOpen.value = false
      showFeedback('Reporte recibido. El cuerpo editorial revisará la observación.')
      return true
    } catch (err: unknown) {
      showFeedback(errorMessage(err, 'Error al enviar el reporte.'))
      return false
    } finally {
      isReportSubmitting.value = false
    }
  }

  return {
    isReportOpen,
    reportTargetType,
    reportTargetId,
    reportTargetLabel,
    isReportSubmitting,
    reportFeedbackMsg,
    openReport,
    handleReportSubmit,
  }
}
