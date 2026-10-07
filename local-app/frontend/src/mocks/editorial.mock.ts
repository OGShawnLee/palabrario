export interface QuorumVoteMock {
  editor_name: string;
  editor_role: string;
  country: string;
  flag: string;
  voted_at: string;
  comment?: string;
}

export interface ProposalMock {
  id_proposal: number;
  lemma: string;
  part_of_speech: string;
  region_code: string;
  region_name: string;
  flag_emoji: string;
  cefr_level: string;
  definition: string;
  short_definition: string;
  cultural_note: string;
  example: string;
  author: string;
  submitted_at: string;
  status: "DRAFT" | "APPROVED" | "REJECTED";
  votes: QuorumVoteMock[];
}

export const MOCK_PROPOSALS: ProposalMock[] = [
  {
    id_proposal: 101,
    lemma: "Zarpar",
    part_of_speech: "verbo intransitivo",
    region_code: "CL",
    region_name: "Chile / Cono Sur",
    flag_emoji: "🇨🇱",
    cefr_level: "B2",
    definition: "Cometer un despropósito, excederse en una actitud o actuar con extrema audacia imprevista.",
    short_definition: "Pasarse del límite o cometer un despropósito.",
    cultural_note: "Modismo contemporáneo derivado del lunfardo portuario chileno y argentino.",
    example: "Te zarpaste con ese comentario, tienes que pedir una disculpa formal.",
    author: "Prof. Alejandro Gómez (Chile)",
    submitted_at: "Hace 2 días",
    status: "DRAFT",
    votes: [
      {
        editor_name: "Prof. Martínez",
        editor_role: "Docente Lingüística",
        country: "Chile",
        flag: "🇨🇱",
        voted_at: "Ayer a las 14:20",
        comment: "Excelente contextualización de la evolución sintáctica.",
      },
      {
        editor_name: "Dra. Rivas",
        editor_role: "Curadora Senior",
        country: "España",
        flag: "🇪🇸",
        voted_at: "Hoy a las 09:15",
        comment: "Cumple con el estándar lexicográfico y claridad pedagógica.",
      },
    ],
  },
  {
    id_proposal: 102,
    lemma: "Chinear",
    part_of_speech: "verbo transitivo",
    region_code: "CR",
    region_name: "Costa Rica / Centroamérica",
    flag_emoji: "🇨🇷",
    cefr_level: "B1",
    definition: "Cuidar, mimar o consentir con afecto excesivo, especialmente a niños o personas queridas.",
    short_definition: "Mimar o consentir con mucho cariño.",
    cultural_note: "Uso arraigado en la identidad costarricense y la hospitalidad pura vida.",
    example: "A mi hermanito le encanta que lo chineen cuando está enfermo.",
    author: "Lic. María Elena Chaves (Costa Rica)",
    submitted_at: "Hace 4 horas",
    status: "DRAFT",
    votes: [
      {
        editor_name: "Mtro. Fonseca",
        editor_role: "Dialectólogo",
        country: "Panamá",
        flag: "🇵🇦",
        voted_at: "Hace 1 hora",
        comment: "Muy fiel al uso pragmático en el istmo.",
      },
    ],
  },
  {
    id_proposal: 103,
    lemma: "Chamba",
    part_of_speech: "sustantivo femenino",
    region_code: "MX",
    region_name: "México",
    flag_emoji: "🇲🇽",
    cefr_level: "A2",
    definition: "Ocupación laboral o puesto de trabajo remunerado.",
    short_definition: "Trabajo o empleo cotidiano.",
    cultural_note: "Término adoptado en el contexto migratorio bracero de los años 40.",
    example: "Salgo de la chamba y nos vemos en la plaza.",
    author: "Dr. Octavio Paz Jr. (México)",
    submitted_at: "Hace 5 días",
    status: "APPROVED",
    votes: [
      { editor_name: "Prof. Ramírez", editor_role: "Editor", country: "México", flag: "🇲🇽", voted_at: "Hace 4 días" },
      { editor_name: "Dra. Rivas", editor_role: "Curadora", country: "España", flag: "🇪🇸", voted_at: "Hace 3 días" },
      { editor_name: "Prof. Soto", editor_role: "Editor", country: "Colombia", flag: "🇨🇴", voted_at: "Hace 2 días" },
    ],
  },
];
