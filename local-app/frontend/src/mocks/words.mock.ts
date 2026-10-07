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
  example?: string;
}

export interface HeteronymMock {
  region_code: string;
  region_name: string;
  flag_emoji: string;
  term: string;
}

export interface WordMock {
  id_word: number;
  lemma: string;
  normalized_lemma: string;
  part_of_speech: string;
  etymology: string;
  friendly_blurb: string;
  origin_flag?: string;
  phonetic?: string;
  heteronyms?: HeteronymMock[];
  senses: SenseMock[];
}

export const MOCK_WORDS: WordMock[] = [
  {
    id_word: 1,
    lemma: "popote",
    normalized_lemma: "popote",
    part_of_speech: "sustantivo masculino",
    phonetic: "/po'po.te/",
    etymology: "Del náhuatl popotl (tallo hueco o paja seca).",
    friendly_blurb:
      "Si vas a una taquería en México y pides refresco, te darán un popote. Pero si cruzas fronteras, prepárate: en Colombia es 'pitillo', en Argentina 'sorbete' o 'pajita', y en Chile 'bombilla'. Una misma pajita, ¡media docena de identidades!",
    origin_flag: "🇲🇽",
    heteronyms: [
      { region_code: "CO", region_name: "Colombia", flag_emoji: "🇨🇴", term: "pitillo" },
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "sorbete" },
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "pajita" },
      { region_code: "CL", region_name: "Chile", flag_emoji: "🇨🇱", term: "bombilla" },
    ],
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
        example: "—¿Me regalas un popote para el agua de horchata, por favor?",
        cultural_note: "En años recientes se ha impulsado el popote de aguacate o bambú para reducir plástico.",
      },
    ],
  },
  {
    id_word: 2,
    lemma: "apapachar",
    normalized_lemma: "apapachar",
    part_of_speech: "verbo transitivo",
    phonetic: "/a.pa.pa't͡ʃaɾ/",
    etymology: "Del náhuatl papatzoa (ablandar algo con caricias o masajes).",
    friendly_blurb:
      "Catalogada muchas veces como la palabra más hermosa del español mexicano. Los lingüistas la definen como 'acariciar con el alma'. Cuando alguien tuvo un mal día, no le des consejos: dale un buen apapacho.",
    origin_flag: "🇲🇽",
    heteronyms: [
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "mimar" },
      { region_code: "CO", region_name: "Colombia", flag_emoji: "🇨🇴", term: "consentir" },
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "acariciar" },
    ],
    senses: [
      {
        id_sense: 102,
        short_definition: "Acariciar o consentir con mucho afecto y ternura profunda.",
        definition: "Dar muestras de cariño especialmente cálidas, brindando confort emocional y cercanía física protectora.",
        region_code: "MX",
        region_name: "México",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "B1",
        example: "Mi abuela siempre me apapachaba preparándome un atole caliente cuando llovía.",
        cultural_note: "Traspasó fronteras hacia Centroamérica y Colombia por su resonancia poética.",
      },
    ],
  },
  {
    id_word: 3,
    lemma: "guagua",
    normalized_lemma: "guagua",
    part_of_speech: "sustantivo",
    phonetic: "/'gwa.gwa/",
    etymology: "Poligenético: del quechua wawa (bebé) en Sudamérica; origen caribeño para el transporte público.",
    friendly_blurb:
      "Una de las palabras trampa más fascinantes del idioma. Si estás en La Habana o en Tenerife y te dicen 'corre a la guagua', súbete al camión. Pero si estás en Santiago de Chile o Quito, ¡guagua es un tierno bebé de brazos!",
    origin_flag: "🌐",
    heteronyms: [
      { region_code: "MX", region_name: "México", flag_emoji: "🇲🇽", term: "camión / bebé" },
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "colectivo / nene" },
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "autobús / crío" },
    ],
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
        example: "Esperamos media hora en la parada pero la guagua venía completamente repleta.",
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
        example: "Abrigaron a la guagüita con una mantita de lana para sacarla al jardín.",
      },
    ],
  },
  {
    id_word: 4,
    lemma: "chamaco",
    normalized_lemma: "chamaco",
    part_of_speech: "sustantivo",
    phonetic: "/t͡ʃa'ma.ko/",
    etymology: "Del maya chamak (niño, muchacho).",
    friendly_blurb:
      "Manera entrañable y coloquial de llamar a los niños y adolescentes en México y Centroamérica. Transmite picardía infantil, energía desbordante y cariño popular.",
    origin_flag: "🇲🇽",
    heteronyms: [
      { region_code: "CO", region_name: "Colombia", flag_emoji: "🇨🇴", term: "pelao" },
      { region_code: "CL", region_name: "Chile", flag_emoji: "🇨🇱", term: "cabro" },
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "chaval" },
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "pibe" },
      { region_code: "VE", region_name: "Venezuela", flag_emoji: "🇻🇪", term: "chamo" },
    ],
    senses: [
      {
        id_sense: 105,
        short_definition: "Niño, infante o muchacho joven.",
        definition: "Persona que está en la etapa de la niñez o adolescencia.",
        region_code: "MX",
        region_name: "México",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "A2",
        example: "Los chamacos del callejón salieron a jugar una cascarita bajo la lluvia.",
        cultural_note: "Muy común en el habla cotidiana urbana y de provincia.",
      },
    ],
  },
  {
    id_word: 5,
    lemma: "coger",
    normalized_lemma: "coger",
    part_of_speech: "verbo transitivo e intransitivo",
    phonetic: "/ko'xeɾ/",
    etymology: "Del latín colligĕre (juntar, recoger, tomar).",
    friendly_blurb:
      "El 'falso amigo' supremo del mundo hispanohablante. Mientras en Madrid tomas el autobús diciendo 'voy a coger el bus', en México y Buenos Aires ese verbo es tabú y tiene una fuerte carga sexual explícita. ¡En América di siempre 'tomar' o 'agarrar'!",
    origin_flag: "🇪🇸",
    heteronyms: [
      { region_code: "MX", region_name: "México", flag_emoji: "🇲🇽", term: "tomar / agarrar" },
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "agarrar" },
    ],
    senses: [
      {
        id_sense: 106,
        short_definition: "Agarrar, asir, tomar algo con la mano o subirse a un transporte.",
        definition: "Tomar o asir una cosa física o abordar un medio de transporte colectivo.",
        region_code: "ES",
        region_name: "España",
        flag_emoji: "🇪🇸",
        is_primary: true,
        cefr_level: "A1",
        example: "Voy a coger el metro de la línea 1 para llegar a tiempo a Sol.",
      },
      {
        id_sense: 107,
        short_definition: "Practicar el coito / relación sexual (Tabú coloquial).",
        definition: "Término vulgar de uso informal para denotar el acto sexual.",
        region_code: "MX",
        region_name: "México / Río de la Plata",
        flag_emoji: "🇲🇽",
        is_primary: false,
        is_taboo: true,
        cultural_note: "En México y Argentina se considera lenguaje vulgar o inapropiado en contextos formales.",
        cefr_level: "B2",
        example: "(Uso coloquial vulgar no recomendado en entornos formales).",
      },
    ],
  },
  {
    id_word: 6,
    lemma: "tianguis",
    normalized_lemma: "tianguis",
    part_of_speech: "sustantivo masculino",
    phonetic: "/'tjaŋ.gis/",
    etymology: "Del náhuatl tiyānquiztli (mercado público).",
    friendly_blurb:
      "El corazón de los barrios mexicanos. No es un supermercado ni un centro comercial: es un laberinto colorido de lonas rosas y amarillas, marchantes gritando ofertas y puestos de comida callejera que se instalan un día a la semana.",
    origin_flag: "🇲🇽",
    heteronyms: [
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "mercadillo" },
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "feria franca" },
      { region_code: "CL", region_name: "Chile", flag_emoji: "🇨🇱", term: "feria libre" },
    ],
    senses: [
      {
        id_sense: 108,
        short_definition: "Mercado ambulante tradicional instalado periódicamente al aire libre.",
        definition: "Bazar o mercado tradicional sobre ruedas que se instala temporalmente en calles de ciudades y pueblos de México.",
        region_code: "MX",
        region_name: "México",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "B1",
        example: "Los domingos vamos al tianguis a comer barbacoa y comprar fruta fresca de temporada.",
      },
    ],
  },
  {
    id_word: 7,
    lemma: "chamba",
    normalized_lemma: "chamba",
    part_of_speech: "sustantivo femenino",
    phonetic: "/'t͡ʃam.ba/",
    etymology: "Probablemente de 'chamber' (cámara de comercio) en el contexto de los braceros mexicanos en EE. UU.",
    friendly_blurb:
      "Término informal por excelencia para referirse al empleo o al trabajo arduo. Si un amigo te dice 'ando en la chamba', significa que está ocupado ganándose el pan.",
    origin_flag: "🇲🇽",
    heteronyms: [
      { region_code: "AR", region_name: "Argentina", flag_emoji: "🇦🇷", term: "laburo" },
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "curro" },
      { region_code: "CL", region_name: "Chile", flag_emoji: "🇨🇱", term: "pega" },
      { region_code: "CO", region_name: "Colombia", flag_emoji: "🇨🇴", term: "camello" },
    ],
    senses: [
      {
        id_sense: 109,
        short_definition: "Trabajo, empleo, ocupación o labor cotidiana.",
        definition: "Puesto laboral o tarea remunerada que desempeña una persona para su sustento.",
        region_code: "MX",
        region_name: "México / Centroamérica",
        flag_emoji: "🇲🇽",
        is_primary: true,
        cefr_level: "A2",
        example: "Termino la chamba a las seis de la tarde y luego nos vemos para el café.",
      },
    ],
  },
  {
    id_word: 8,
    lemma: "zarpar",
    normalized_lemma: "zarpar",
    part_of_speech: "verbo intransitivo",
    phonetic: "/saɾ'paɾ/",
    etymology: "Metáfora náutica: desanclar la embarcación; en lunfardo: pasarse de la raya.",
    friendly_blurb:
      "En el lenguaje marino tradicional significa levar anclas y salir al mar. Pero en el cono sur (Argentina, Uruguay, Chile) es un modismo juvenil vital: 'zarpar' o 'zarparse' significa pasarse del límite o hacer algo increíblemente osado.",
    origin_flag: "🇦🇷",
    heteronyms: [
      { region_code: "MX", region_name: "México", flag_emoji: "🇲🇽", term: "pasarse de lanza" },
      { region_code: "ES", region_name: "España", flag_emoji: "🇪🇸", term: "pasarse de la raya" },
    ],
    senses: [
      {
        id_sense: 110,
        short_definition: "Cometer un despropósito, excederse en una actitud o actuar de forma osada.",
        definition: "Sobrepasar los límites de lo prudente, educado o convencional en un comportamiento social.",
        region_code: "AR",
        region_name: "Argentina / Cono Sur",
        flag_emoji: "🇦🇷",
        is_primary: true,
        cefr_level: "B2",
        example: "Te zarpaste con ese chiste frente a los profesores, tienes que pedir una disculpa.",
        cultural_note: "Evolución contemporánea del lunfardo de Buenos Aires.",
      },
    ],
  },
];
