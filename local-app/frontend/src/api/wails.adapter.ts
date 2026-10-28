import * as WailsAppBridge from '../../bindings/palabrario/local-app/app';
import type { IPalabrarioAPI } from './client';

/**
 * Adaptador de Transporte de Escritorio (Wails v3 IPC).
 * Invoca directamente los bindings nativos autogenerados desde local-app/app.go.
 */
export const WailsDesktopAdapter: IPalabrarioAPI = WailsAppBridge;
