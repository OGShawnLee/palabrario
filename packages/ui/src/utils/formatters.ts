/**
 * Helpers y diccionarios canónicos de localización al español para microcopy léxico.
 * Evita la fuga de enums o valores crudos en inglés (male, everyday, noun, etc.).
 */

export const GENDER_LABELS: Record<string, string> = {
  male: 'masculino',
  MALE: 'masculino',
  female: 'femenino',
  FEMALE: 'femenino',
  neutral: 'neutro',
  NEUTRAL: 'neutro',
  ambiguous: 'ambiguo',
  AMBIGUOUS: 'ambiguo',
}

export function formatGender(gender?: string | null): string {
  if (!gender) return ''
  return GENDER_LABELS[gender] || gender.toLowerCase()
}

export const FREQUENCY_LABELS: Record<string, string> = {
  EVERYDAY: 'Uso cotidiano',
  everyday: 'Uso cotidiano',
  COMMON: 'Común',
  common: 'Común',
  RARE: 'Poco frecuente',
  rare: 'Poco frecuente',
  SLANG: 'Jerga / Argot',
  slang: 'Jerga / Argot',
  ARCHAIC: 'Arcaico',
  archaic: 'Arcaico',
}

export function formatFrequency(frequency?: string | null): string {
  if (!frequency) return ''
  return FREQUENCY_LABELS[frequency] || frequency
}

export const PART_OF_SPEECH_LABELS: Record<string, string> = {
  noun: 'sustantivo',
  sustantivo: 'sustantivo',
  verb: 'verbo',
  verbo: 'verbo',
  adjective: 'adjetivo',
  adjetivo: 'adjetivo',
  adverb: 'adverbio',
  adverbio: 'adverbio',
  idiom: 'modismo',
  modismo: 'modismo',
  interjection: 'interjección',
  interjección: 'interjección',
  phrase: 'frase / locución',
  frase: 'frase / locución',
}

export function formatPartOfSpeech(pos?: string | null): string {
  if (!pos) return ''
  const normalized = pos.trim().toLowerCase()
  return PART_OF_SPEECH_LABELS[normalized] || pos
}
