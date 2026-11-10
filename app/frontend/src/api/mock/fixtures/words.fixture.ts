import {
  Example,
  HeteronymItem,
  Region,
  SenseDetail,
  WordDetail,
  WordGender,
} from "@bindings/palabrario/pkg/models/index.js"

// Canonical dialect codes: 'MX', 'ES', 'AR', 'CO', 'CL', 'CU', 'CAM', 'PAN'
// UI rendering rule: never render flag_emoji or raw codes directly in UI components.

const TIMESTAMP = "2026-01-01T00:00:00.000Z"

function region(
  id: string,
  code: string,
  name: string,
  parentZone: string,
  countryCode: string | null = null,
  flagEmoji = "",
): Region {
  return new Region({
    id_region: id,
    code,
    name,
    parent_zone: parentZone,
    country_code: countryCode,
    flag_emoji: flagEmoji,
  })
}

function heteronym(
  regionCode: string,
  regionName: string,
  term: string,
  countryCode = "",
): HeteronymItem {
  return new HeteronymItem({
    region_code: regionCode,
    region_name: regionName,
    country_code: countryCode || (regionCode.length === 2 ? regionCode : ""),
    flag_emoji: "",
    term,
  })
}

export const mockRegions: Region[] = [
  region("a1111111-0000-0000-0000-000000000001", "MX", "México", "Norteamérica", "MX"),
  region("a1111111-0000-0000-0000-000000000002", "ES", "España", "Peninsular", "ES"),
  region("a1111111-0000-0000-0000-000000000003", "AR", "Argentina", "Rioplatense", "AR"),
  region("a1111111-0000-0000-0000-000000000004", "CO", "Colombia", "Andina", "CO"),
  region("a1111111-0000-0000-0000-000000000005", "CL", "Chile", "Cono Sur", "CL"),
  region("a1111111-0000-0000-0000-000000000006", "CU", "Cuba", "Caribe", "CU"),
  region("a1111111-0000-0000-0000-000000000007", "CAM", "Centroamérica", "Centroamérica", null),
  region("a1111111-0000-0000-0000-000000000008", "PAN", "Panhispánico", "General", null),
]

export const guaguaWord: WordDetail = new WordDetail({
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
    heteronym("MX", "México", "Camión"),
    heteronym("AR", "Argentina", "Colectivo"),
    heteronym("CL", "Chile", "Micro"),
    heteronym("ES", "España", "Autobús"),
    heteronym("CL", "Chile", "Guagua"),
    heteronym("ES", "España", "Bebé"),
    heteronym("MX", "México", "Bebé"),
    heteronym("AR", "Argentina", "Bebé"),
  ],
  created_at: TIMESTAMP,
  updated_at: TIMESTAMP,
})

const cheWord: WordDetail = new WordDetail({
  id_word: 102,
  lemma: "Che",
  normalized_lemma: "che",
  part_of_speech: "interjección",
  gender: WordGender.GenderNeutral,
  etymology:
    "Probablemente del mapuche che (gente, persona), o del valenciano antiguo; se popularizó en el Río de la Plata como vocativo de confianza.",
  friendly_blurb:
    "La palabra que delata a un rioplatense en cualquier aeropuerto del mundo. Sirve para llamar la atención, pedir algo, mostrar sorpresa o simplemente marcar cercanía con quien conversas.",
  image_url: null,
  origin_region: null,
  senses: [
    new SenseDetail({
      id_sense: 211,
      id_word: 102,
      region_code: "AR",
      region_name: "Argentina",
      country_code: "AR",
      flag_emoji: "",
      short_definition: "Interjección para llamar la atención de alguien",
      definition:
        "Interjección coloquial que se usa para llamar a una persona, expresar sorpresa o dar énfasis a una frase en un trato cercano.",
      is_primary: true,
      is_taboo_or_vulgar: false,
      cultural_note:
        "Tan característica es esta voz que dio nombre al apodo más célebre de la región. En una charla entre amigos puede aparecer varias veces por frase.",
      image_url: null,
      cefr_level: "A2",
      usage_frequency: "EVERYDAY",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 311,
          id_sense: 211,
          sentence: "Che, ¿me pasás el mate?",
          source: "Habla cotidiana · Argentina",
        }),
        new Example({
          id_example: 312,
          id_sense: 211,
          sentence: "¡Che, no puede ser que llegues tarde otra vez!",
          source: "Habla cotidiana · Argentina",
        }),
      ],
    }),
  ],
  heteronyms: [
    heteronym("MX", "México", "Oye"),
    heteronym("ES", "España", "Tío"),
    heteronym("CO", "Colombia", "Parce"),
  ],
  created_at: TIMESTAMP,
  updated_at: TIMESTAMP,
})

const bacanWord: WordDetail = new WordDetail({
  id_word: 103,
  lemma: "Bacán",
  normalized_lemma: "bacan",
  part_of_speech: "adjetivo",
  gender: WordGender.GenderAmbiguous,
  etymology:
    "De origen incierto; en Chile y Perú se asocia con el genovés baccan (patrón, jefe), usado por los inmigrantes italianos del siglo XIX.",
  friendly_blurb:
    "Cuando algo es genial, estupendo o simplemente muy bueno, en los Andes lo llaman bacán. Una palabra de aprobación que suena a buen ánimo.",
  image_url: null,
  origin_region: null,
  senses: [
    new SenseDetail({
      id_sense: 221,
      id_word: 103,
      region_code: "CL",
      region_name: "Chile",
      country_code: "CL",
      flag_emoji: "",
      short_definition: "Excelente, muy bueno",
      definition:
        "Dicho de una cosa o de una situación: que resulta agradable, estupenda o de gran calidad.",
      is_primary: true,
      is_taboo_or_vulgar: false,
      cultural_note:
        "Se usa sobre todo entre jóvenes, aunque hoy se entiende en todas las generaciones. Funciona como elogio rápido ante un plan, una canción o una comida.",
      image_url: null,
      cefr_level: "B1",
      usage_frequency: "COMMON",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 321,
          id_sense: 221,
          sentence: "Estuvo bacán el concierto de anoche.",
          source: "Habla juvenil · Chile",
        }),
      ],
    }),
    new SenseDetail({
      id_sense: 222,
      id_word: 103,
      region_code: "CO",
      region_name: "Colombia",
      country_code: "CO",
      flag_emoji: "",
      short_definition: "Persona elegante o de buena posición",
      definition:
        "Persona que viste con elegancia o que goza de una posición económica holgada.",
      is_primary: false,
      is_taboo_or_vulgar: false,
      cultural_note: null,
      image_url: null,
      cefr_level: "B2",
      usage_frequency: "COMMON",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 322,
          id_sense: 222,
          sentence: "Llegó hecho un bacán, con traje y zapatos nuevos.",
          source: "Habla popular · Colombia",
        }),
      ],
    }),
  ],
  heteronyms: [
    heteronym("MX", "México", "Padre"),
    heteronym("ES", "España", "Guay"),
    heteronym("AR", "Argentina", "Copado"),
    heteronym("CO", "Colombia", "Chévere"),
  ],
  created_at: TIMESTAMP,
  updated_at: TIMESTAMP,
})

const chuchoWord: WordDetail = new WordDetail({
  id_word: 104,
  lemma: "Chucho",
  normalized_lemma: "chucho",
  part_of_speech: "sustantivo masculino",
  gender: WordGender.GenderMale,
  etymology:
    "De origen onomatopéyico, a partir de la voz con la que se ahuyenta o llama a los perros; en México también es hipocorístico del nombre Jesús.",
  friendly_blurb:
    "Una palabra con vidas paralelas: para unos es el perro sin raza del vecindario, para otros un escalofrío repentino, y en México puede ser simplemente un tocayo llamado Jesús.",
  image_url: null,
  origin_region: null,
  senses: [
    new SenseDetail({
      id_sense: 231,
      id_word: 104,
      region_code: "CO",
      region_name: "Colombia",
      country_code: "CO",
      flag_emoji: "",
      short_definition: "Perro callejero o sin raza definida",
      definition:
        "Perro mestizo, generalmente de tamaño mediano o pequeño, que no pertenece a una raza concreta.",
      is_primary: true,
      is_taboo_or_vulgar: false,
      cultural_note:
        "Aunque a veces se usa con cierto desdén, también es un apelativo cariñoso para el perro de la casa que nadie sabría cómo reemplazar.",
      image_url: null,
      cefr_level: "B1",
      usage_frequency: "COMMON",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 331,
          id_sense: 231,
          sentence: "El chucho de la esquina siempre me acompaña hasta el colegio.",
          source: "Habla cotidiana · Colombia",
        }),
      ],
    }),
    new SenseDetail({
      id_sense: 232,
      id_word: 104,
      region_code: "CL",
      region_name: "Chile",
      country_code: "CL",
      flag_emoji: "",
      short_definition: "Escalofrío o temblor por frío o miedo",
      definition:
        "Sensación repentina de frío o de miedo que recorre el cuerpo y provoca un estremecimiento.",
      is_primary: false,
      is_taboo_or_vulgar: false,
      cultural_note:
        "Casi siempre aparece en plural y con el verbo dar: dar chuchos, es decir, temblar de frío.",
      image_url: null,
      cefr_level: "B2",
      usage_frequency: "COMMON",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 332,
          id_sense: 232,
          sentence: "Con este viento me dan chuchos de solo salir a la calle.",
          source: "Habla cotidiana · Chile",
        }),
      ],
    }),
  ],
  heteronyms: [
    heteronym("MX", "México", "Perro"),
    heteronym("ES", "España", "Chucho"),
    heteronym("AR", "Argentina", "Perro callejero"),
  ],
  created_at: TIMESTAMP,
  updated_at: TIMESTAMP,
})

const parceWord: WordDetail = new WordDetail({
  id_word: 105,
  lemma: "Parce",
  normalized_lemma: "parce",
  part_of_speech: "sustantivo común",
  gender: WordGender.GenderAmbiguous,
  etymology:
    "Apócope de parcero, voz paisa de Antioquia que probablemente deriva de parcial o de la expresión paisa parcero, compañero de andanzas.",
  friendly_blurb:
    "El saludo por excelencia de Medellín y de buena parte de Colombia. Con un parce nunca estás solo: es amigo, cómplice y compañero de plan, todo en una sola sílaba.",
  image_url: null,
  origin_region: null,
  senses: [
    new SenseDetail({
      id_sense: 241,
      id_word: 105,
      region_code: "CO",
      region_name: "Colombia",
      country_code: "CO",
      flag_emoji: "",
      short_definition: "Amigo cercano o compañero",
      definition:
        "Amigo de confianza con quien se comparte la vida cotidiana; también se usa como vocativo para dirigirse a alguien con cercanía.",
      is_primary: true,
      is_taboo_or_vulgar: false,
      cultural_note:
        "Se escucha con especial fuerza en Antioquia y el Eje Cafetero. Decirlo a un desconocido puede sonar amistoso o atrevido, según el tono.",
      image_url: null,
      cefr_level: "B1",
      usage_frequency: "EVERYDAY",
      state: "APPROVED",
      examples: [
        new Example({
          id_example: 341,
          id_sense: 241,
          sentence: "¡Qué hubo, parce! ¿Para dónde va con tanta prisa?",
          source: "Habla paisa · Colombia",
        }),
        new Example({
          id_example: 342,
          id_sense: 241,
          sentence: "Mi parce me ayudó con la mudanza sin pedir nada a cambio.",
          source: "Habla cotidiana · Colombia",
        }),
      ],
    }),
  ],
  heteronyms: [
    heteronym("MX", "México", "Cuate"),
    heteronym("ES", "España", "Colega"),
    heteronym("AR", "Argentina", "Pibe"),
    heteronym("CL", "Chile", "Compadre"),
  ],
  created_at: TIMESTAMP,
  updated_at: TIMESTAMP,
})

export const mockWords: WordDetail[] = [
  guaguaWord,
  cheWord,
  bacanWord,
  chuchoWord,
  parceWord,
]

/** Id of the fixture used as the default word of the day. */
export const MOCK_WORD_OF_THE_DAY_ID = guaguaWord.id_word
