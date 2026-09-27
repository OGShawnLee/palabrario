export interface SenseMock {
  id_sense: number;
  short_definition: string;
  definition: string;
  region_code: string;
  region_name: string;
  flag_emoji: string;
  is_primary: boolean;
  is_taboo?: boolean;
  cultural_note?: string;
  cefr_level: string;
}

export interface WordMock {
  id_word: number;
  lemma: string;
  normalized_lemma: string;
  part_of_speech: string;
  etymology: string;
  friendly_blurb: string;
  origin_flag?: string;
  senses: SenseMock[];
}

export const MOCK_WORDS: WordMock[] = [
  {
    id_word: 1,
    lemma: "popote",
    normalized_lemma: "popote",
    part_of_speech: "sustantivo masculino",
    etymology: "Del náhuatl popotl (tallo hueco o paja seca).",
    friendly_blurb:
      "Si vas a una taquería en México y pides refresco, te darán un popote. Pero si cruzas fronteras, prepárate: en Colombia es 'pitillo', en Argentina 'pajita' o 'sorbete', y en Chile 'bombilla'. Una misma pajita, ¡media docena de identidades!",
    origin_flag: "🇲🇽",
    senses: [
      {
        id_sense: 101,
        short_definition: "Tubo hueco de plástico o papel para sorber bebidas.",
        definition: "Cilindro hueco y delgado utilizado comúnmente para extraer líquidos de un recipiente hacia la boca.",
        region_code: "MX",
        region_name: "México",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "A2",
      },
    ],
  },
  {
    id_word: 2,
    lemma: "apapachar",
    normalized_lemma: "apapachar",
    part_of_speech: "verbo transitivo",
    etymology: "Del náhuatl papatzoa (ablandar algo con caricias o masajes).",
    friendly_blurb:
      "Catalogada muchas veces como la palabra más hermosa del español mexicano. Los lingüistas la definen como 'acariciar con el alma'. Cuando alguien tuvo un mal día, no le des consejos: dale un buen apapacho.",
    origin_flag: "🇲🇽",
    senses: [
      {
        id_sense: 102,
        short_definition: "Acariciar o consentir con mucho afecto y ternura.",
        definition: "Dar muestras de cariño especialmente cálidas, brindando confort emocional y cercanía física.",
        region_code: "MX",
        region_name: "México",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "B1",
      },
    ],
  },
  {
    id_word: 3,
    lemma: "guagua",
    normalized_lemma: "guagua",
    part_of_speech: "sustantivo",
    etymology: "Poligenético: del quechua wawa (bebé) en Sudamérica; origen caribeño para el transporte público.",
    friendly_blurb:
      "Una de las palabras trampa más fascinantes del idioma. Si estás en La Habana o en Tenerife y te dicen 'corre a la guagua', súbete al camión. Pero si estás en Santiago de Chile o Quito, ¡guagua es un tierno bebé de brazos!",
    origin_flag: "🌐",
    senses: [
      {
        id_sense: 103,
        short_definition: "Autobús o vehículo de transporte colectivo público.",
        definition: "Vehículo automotor de gran capacidad destinado al transporte regular de pasajeros por vía urbana o interurbana.",
        region_code: "CU",
        region_name: "Cuba / Caribe / Canarias",
        flag_emoji: "🇨🇺",
        is_primary: true,
        cultural_note: "En las Antillas hispanas y las Islas Canarias es el término estándar diario.",
        cefr_level: "B1",
      },
      {
        id_sense: 104,
        short_definition: "Niño pequeño, infante o criatura recién nacida.",
        definition: "Bebé de muy corta edad o que aún no camina.",
        region_code: "CL",
        region_name: "Chile / Región Andina",
        flag_emoji: "🇨🇱",
        is_primary: false,
        cultural_note: "Proviene de la lengua quechua; uso generalizado en Chile, Bolivia, Ecuador y Perú.",
        cefr_level: "B1",
      },
    ],
  },
  {
    id_word: 4,
    lemma: "coger",
    normalized_lemma: "coger",
    part_of_speech: "verbo transitivo e intransitivo",
    etymology: "Del latín colligĕre (juntar, recoger, tomar).",
    friendly_blurb:
      "El 'falso amigo' supremo del mundo hispanohablante. Mientras en Madrid tomas el autobús diciendo 'voy a coger el bus', en México y Buenos Aires ese verbo es tabú y tiene una fuerte carga sexual explícita. ¡En América di siempre 'tomar' o 'agarrar'!",
    origin_flag: "🇪🇸",
    senses: [
      {
        id_sense: 105,
        short_definition: "Agarrar, asir, tomar algo con la mano o subirse a un transporte.",
        definition: "Tomar o asir una cosa física o abordar un medio de transporte.",
        region_code: "ES",
        region_name: "España",
        flag_emoji: "🇪🇸",
        is_primary: true,
        cefr_level: "A1",
      },
      {
        id_sense: 106,
        short_definition: "Practicar el coito / relación sexual (Tabú coloquial).",
        definition: "Término vulgar de uso informal para denotar el acto sexual.",
        region_code: "MX",
        region_name: "México / Río de la Plata",
        flag_emoji: "🇲🇽",
        is_primary: false,
        is_taboo: true,
        cultural_note: "En México y Argentina se considera lenguaje vulgar o inapropiado en contextos formales.",
        cefr_level: "B2",
      },
    ],
  },
  {
    id_word: 5,
    lemma: "tianguis",
    normalized_lemma: "tianguis",
    part_of_speech: "sustantivo masculino",
    etymology: "Del náhuatl tiyānquiztli (mercado público).",
    friendly_blurb:
      "El corazón de los barrios mexicanos. No es un supermercado ni un centro comercial: es un laberinto colorido de lonas rosas y amarillas, marchantes gritando ofertas y puestos de comida callejera que se instalan un día a la semana.",
    origin_flag: "🇲🇽",
    senses: [
      {
        id_sense: 107,
        short_definition: "Mercado ambulante tradicional instalado periódicamente al aire libre.",
        definition: "Bazar o mercado tradicional sobre ruedas que se instala temporalmente en calles de ciudades y pueblos de México.",
        region_code: "MX",
        region_name: "México",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "B1",
      },
    ],
  },
  {
    id_word: 6,
    lemma: "chamba",
    normalized_lemma: "chamba",
    part_of_speech: "sustantivo femenino",
    etymology: "Probablemente de 'chamber' (cámara de comercio) en el contexto de los braceros mexicanos en EE. UU.",
    friendly_blurb:
      "Término informal por excelencia para referirse al empleo o al trabajo arduo. Si un amigo te dice 'ando en la chamba', significa que está ocupado ganándose el pan.",
    origin_flag: "🇲🇽",
    senses: [
      {
        id_sense: 108,
        short_definition: "Trabajo, empleo, ocupación o labor cotidiana.",
        definition: "Puesto laboral o tarea remunerada que desempeña una persona.",
        region_code: "MX",
        region_name: "México / Centroamérica",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "A2",
      },
    ],
  },
];
