import type { IPalabrarioAPI } from "@/api/contract"
import { WailsAuthAdapter } from "./auth.wails"
import { WailsDictionaryAdapter } from "./dictionary.wails"
import { WailsStudioAdapter } from "./studio.wails"

export class WailsAdapter implements IPalabrarioAPI {
  readonly auth = new WailsAuthAdapter()
  readonly dictionary = new WailsDictionaryAdapter()
  readonly studio = new WailsStudioAdapter()
}
