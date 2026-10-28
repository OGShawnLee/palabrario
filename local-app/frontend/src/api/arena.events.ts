import { Events } from '@wailsio/runtime';
import type * as models$0 from '../../bindings/palabrario/pkg/models/models.js';
import { isWailsRuntime } from './client';

export type ArenaMessageHandler = (message: models$0.WSEvent) => void;

class ArenaEventBusClass {
  private browserListeners: Set<ArenaMessageHandler> = new Set();

  /**
   * Suscribe un listener a tramas WSEvent unificando Wails (Events.On) y Web/Capacitor.
   * @param handler Función que recibirá cada trama entrante WSEvent.
   * @returns Función de desuscripción () => void.
   */
  public onMessage(handler: ArenaMessageHandler): () => void {
    if (isWailsRuntime()) {
      return Events.On('arena:ws_message', (ev: any) => {
        if (ev && ev.data) {
          handler(ev.data as models$0.WSEvent);
        }
      });
    }

    this.browserListeners.add(handler);
    return () => {
      this.browserListeners.delete(handler);
    };
  }

  /**
   * Método interno para que CapacitorMobileAdapter difunda tramas entrantes en modo Web/Capacitor.
   */
  public dispatchBrowserMessage(message: models$0.WSEvent): void {
    this.browserListeners.forEach((listener) => {
      try {
        listener(message);
      } catch (err) {
        console.error('Error procesando evento en ArenaEventBus:', err);
      }
    });
  }
}

export const ArenaEventBus = new ArenaEventBusClass();
