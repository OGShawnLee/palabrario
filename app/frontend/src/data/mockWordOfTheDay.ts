import {
  WordDetail,
  SenseDetail,
  Example,
  HeteronymItem,
  WordGender,
} from "@bindings/palabrario/pkg/models/index.js"

export interface HeteronymGroup {
  label: string
  items: HeteronymItem[]
}

export const mockHeteronymGroups: HeteronymGroup[] = [
  {
    label: "Para transporte de pasajeros",
    items: [
      new HeteronymItem({
        region_code: "MX",
        region_name: "Mexico",
        country_code: "MX",
        flag_emoji: "",
        term: "Camión",
      }),
      new HeteronymItem({
        region_code: "AR",
        region_name: "Argentina",
        country_code: "AR",
        flag_emoji: "",
        term: "Colectivo",
      }),
      new HeteronymItem({
        region_code: "CL",
        region_name: "Chile",
        country_code: "CL",
        flag_emoji: "",
        term: "Micro",
      }),
      new HeteronymItem({
        region_code: "ES",
        region_name: "Spain",
        country_code: "ES",
        flag_emoji: "",
        term: "Autobús",
      }),
    ],
  },
  {
    label: "Para infancia",
    items: [
      new HeteronymItem({
        region_code: "CL",
        region_name: "Chile",
        country_code: "CL",
        flag_emoji: "",
        term: "Guagua",
      }),
      new HeteronymItem({
        region_code: "ES",
        region_name: "Spain",
        country_code: "ES",
        flag_emoji: "",
        term: "Bebé",
      }),
      new HeteronymItem({
        region_code: "MX",
        region_name: "Mexico",
        country_code: "MX",
        flag_emoji: "",
        term: "Bebé",
      }),
      new HeteronymItem({
        region_code: "AR",
        region_name: "Argentina",
        country_code: "AR",
        flag_emoji: "",
        term: "Bebé",
      }),
    ],
  },
]

export const mockWordOfTheDay: WordDetail = new WordDetail({
  id_word: 101,
  lemma: "Guagua",
  normalized_lemma: "guagua",
  part_of_speech: "sustantivo femenino",
  gender: WordGender.GenderFemale,
  etymology:
    "Del quechua wawa (bebé, criatura) en la zona andina; en el Caribe, probablemente voz onomatopéyica o adaptada de la empresa de transporte Wa & Wa Co. Inc.",
  friendly_blurb:
    "Uno de los camaleones lingüísticos más fascinantes del español. En el Caribe, una guagua es un autobús ruidoso que te lleva al trabajo; en Chile, Ecuador y Perú, una guagua es una criatura recién nacida. ¡Dila en el país equivocado y podrías terminar cargando un autobús!",
  image_url: null,
  origin_region: null,
  senses: [
    new SenseDetail({
      id_sense: 201,
      id_word: 101,
      region_code: "CU",
      region_name: "Cuba",
      country_code: "CU",
      flag_emoji: "",
      short_definition: "Autobús de transporte público",
      definition:
        "Ómnibus de servicio urbano o interurbano destinado al transporte masivo de pasajeros.",
      is_primary: true,
      is_taboo_or_vulgar: false,
      cultural_note:
        "En Cuba, tomar la guagua forma parte de la vida cotidiana: conecta barrios, trabajos y encuentros, y siempre deja espacio para una conversación.",
      image_url: null,
      cefr_level: "B1",
      usage_frequency: "EVERYDAY",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 301,
          id_sense: 201,
          sentence: "Ayer la guagua pasó llena y tuve que caminar.",
          source: "Habla urbana · Cuba",
        }),
      ],
    }),
    new SenseDetail({
      id_sense: 202,
      id_word: 101,
      region_code: "CL",
      region_name: "Chile",
      country_code: "CL",
      flag_emoji: "",
      short_definition: "Bebé o criatura de pecho",
      definition: "Niño de pecho o criatura de muy corta edad.",
      is_primary: false,
      is_taboo_or_vulgar: false,
      cultural_note:
        "En Chile, guagua es una palabra cotidiana y afectuosa para hablar de los más pequeños. Se escucha en casa, en la calle y en las conversaciones de familia.",
      image_url: null,
      cefr_level: "A2",
      usage_frequency: "EVERYDAY",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 302,
          id_sense: 202,
          sentence: "La guagua se quedó dormida en brazos de su mamá.",
          source: "Habla cotidiana · Chile",
        }),
      ],
    }),
  ],
  heteronyms: [
    ...mockHeteronymGroups[0].items,
    ...mockHeteronymGroups[1].items,
  ],
  created_at: "2026-01-01T00:00:00.000Z",
  updated_at: "2026-01-01T00:00:00.000Z",
})
