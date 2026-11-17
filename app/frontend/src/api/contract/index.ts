import type { IArenaAPI } from "./arena.contract"
import type { IAuthAPI } from "./auth.contract"
import type { IDictionaryAPI } from "./dictionary.contract"
import type { IMasteryAPI } from "./mastery.contract"
import type { IStudioAPI } from "./studio.contract"

export type { IArenaAPI } from "./arena.contract"
export type { IAuthAPI } from "./auth.contract"
export type { IDictionaryAPI } from "./dictionary.contract"
export type { IMasteryAPI } from "./mastery.contract"
export type { IStudioAPI } from "./studio.contract"

export interface IPalabrarioAPI {
  readonly arena: IArenaAPI
  readonly auth: IAuthAPI
  readonly dictionary: IDictionaryAPI
  readonly mastery: IMasteryAPI
  readonly studio: IStudioAPI
}
