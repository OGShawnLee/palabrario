import { CancellablePromise as $CancellablePromise } from '@wailsio/runtime';
import { Preferences } from '@capacitor/preferences';
import type * as models$0 from '../../bindings/palabrario/pkg/models/models.js';
import type { IPalabrarioAPI } from './client';
import { ArenaEventBus } from './arena.events';

const CORE_URL = (import.meta as any).env?.VITE_CORE_URL || 'http://localhost:8000';
const IMAGE_URL = (import.meta as any).env?.VITE_IMAGE_URL || 'http://localhost:8002';
const REALTIME_URL = (import.meta as any).env?.VITE_REALTIME_URL || 'ws://localhost:8001';

const TOKEN_KEY = 'palabrario_session_token';
let cachedToken: string | null = null;
let activeWebSocket: WebSocket | null = null;

export async function getSessionToken(): Promise<string | null> {
  if (cachedToken !== null) return cachedToken;
  try {
    const { value } = await Preferences.get({ key: TOKEN_KEY });
    cachedToken = value;
    return value;
  } catch {
    return null;
  }
}

export async function setSessionToken(token: string): Promise<void> {
  cachedToken = token;
  try {
    await Preferences.set({ key: TOKEN_KEY, value: token });
  } catch (e) {
    console.error('Error guardando token en Preferences:', e);
  }
}

export async function clearSessionToken(): Promise<void> {
  cachedToken = '';
  try {
    await Preferences.remove({ key: TOKEN_KEY });
  } catch (e) {
    console.error('Error eliminando token de Preferences:', e);
  }
}

function toCancellable<T>(promise: Promise<T>): $CancellablePromise<T> {
  return new $CancellablePromise<T>((resolve, reject) => {
    promise.then(resolve, reject);
  });
}

async function httpReq<T>(
  method: string,
  url: string,
  body?: unknown,
  customHeaders?: Record<string, string>
): Promise<T> {
  const token = await getSessionToken();
  const headers: Record<string, string> = {
    Accept: 'application/json',
    ...customHeaders,
  };

  if (body !== undefined && !(body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    method,
    headers,
    body: body instanceof FormData ? body : body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    let errMsg = response.statusText;
    try {
      const errJson = await response.json();
      if (errJson && errJson.error) {
        errMsg = errJson.error;
      }
    } catch {
      // ignore
    }
    throw new Error(errMsg);
  }

  if (response.status === 204) {
    return null as T;
  }

  return (await response.json()) as T;
}

export const CapacitorMobileAdapter: IPalabrarioAPI = {
  HealthCheck(): $CancellablePromise<string> {
    return toCancellable(Promise.resolve('ok'));
  },

  SetContext(): $CancellablePromise<void> {
    return toCancellable(Promise.resolve());
  },

  ModelBindingsBridge(): $CancellablePromise<void> {
    return toCancellable(Promise.resolve());
  },

  // ===========================================================================
  // MÓDULO 1: Sesión & Autenticación
  // ===========================================================================

  HasActiveSession(): $CancellablePromise<boolean> {
    return toCancellable((async () => {
      const token = await getSessionToken();
      return !!token;
    })());
  },

  SignIn(req: models$0.SignInRequest): $CancellablePromise<models$0.UserProfile | null> {
    return toCancellable((async () => {
      const resp = await httpReq<models$0.AuthResponse>('POST', `${CORE_URL}/api/v1/auth/sign-in`, req);
      if (resp?.token) {
        await setSessionToken(resp.token);
      }
      return resp ? resp.user : null;
    })());
  },

  SignUp(req: models$0.SignUpRequest): $CancellablePromise<models$0.UserProfile | null> {
    return toCancellable((async () => {
      const resp = await httpReq<models$0.AuthResponse>('POST', `${CORE_URL}/api/v1/auth/sign-up`, req);
      if (resp?.token) {
        await setSessionToken(resp.token);
      }
      return resp ? resp.user : null;
    })());
  },

  SignOut(): $CancellablePromise<void> {
    return toCancellable((async () => {
      await clearSessionToken();
      if (activeWebSocket) {
        activeWebSocket.close();
        activeWebSocket = null;
      }
    })());
  },

  GetMe(): $CancellablePromise<models$0.UserProfile | null> {
    return toCancellable(httpReq<models$0.UserProfile>('GET', `${CORE_URL}/api/v1/auth/me`));
  },

  GetMyProfile(): $CancellablePromise<models$0.UserProfile | null> {
    return this.GetMe();
  },

  // ===========================================================================
  // MÓDULO 2: Identidad & Perfil Social
  // ===========================================================================

  GetPublicProfileByUsername(username: string): $CancellablePromise<models$0.PublicUserProfile | null> {
    return toCancellable(
      httpReq<models$0.PublicUserProfile>('GET', `${CORE_URL}/api/v1/users/u/${encodeURIComponent(username)}`)
    );
  },

  UpdateProfile(req: models$0.UpdateProfileRequest): $CancellablePromise<models$0.UserProfile | null> {
    return toCancellable(httpReq<models$0.UserProfile>('PUT', `${CORE_URL}/api/v1/users/me`, req));
  },

  ChangePassword(req: models$0.ChangePasswordRequest): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('PUT', `${CORE_URL}/api/v1/users/me/password`, req));
  },

  // ===========================================================================
  // MÓDULO 3: Gobernanza Institucional, Roles y Boletos
  // ===========================================================================

  RedeemInvitationTicket(req: models$0.RedeemTicketRequest): $CancellablePromise<models$0.UserProfile | null> {
    return toCancellable(httpReq<models$0.UserProfile>('POST', `${CORE_URL}/api/v1/governance/tickets/redeem`, req));
  },

  CreateEditorRequest(req: models$0.CreateEditorRequestDTO): $CancellablePromise<models$0.EditorRequest | null> {
    return toCancellable(httpReq<models$0.EditorRequest>('POST', `${CORE_URL}/api/v1/governance/editor-requests`, req));
  },

  ListMyEditorRequests(): $CancellablePromise<models$0.EditorRequest[]> {
    return toCancellable((async () => {
      const res = await httpReq<models$0.EditorRequest[]>('GET', `${CORE_URL}/api/v1/governance/editor-requests/me`);
      return res || [];
    })());
  },

  ListPendingEditorRequests(): $CancellablePromise<models$0.EditorRequest[]> {
    return this.ListEditorRequests('PENDING');
  },

  ListEditorRequests(status: string): $CancellablePromise<models$0.EditorRequest[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/governance/editor-requests${status ? `?status=${encodeURIComponent(status)}` : ''}`;
      const res = await httpReq<models$0.EditorRequest[]>('GET', url);
      return res || [];
    })());
  },

  ResolveEditorRequest(id: string, req: models$0.ResolveEditorRequestDTO): $CancellablePromise<models$0.EditorRequest | null> {
    return toCancellable(
      httpReq<models$0.EditorRequest>('PUT', `${CORE_URL}/api/v1/governance/editor-requests/${encodeURIComponent(id)}/resolve`, req)
    );
  },

  GenerateInvitationTickets(req: models$0.CreateTicketRequest): $CancellablePromise<models$0.InvitationTicket | null> {
    return this.CreateInvitationTicket(req);
  },

  CreateInvitationTicket(req: models$0.CreateTicketRequest): $CancellablePromise<models$0.InvitationTicket | null> {
    return toCancellable(httpReq<models$0.InvitationTicket>('POST', `${CORE_URL}/api/v1/governance/tickets`, req));
  },

  ListInvitationTickets(activeOnly: boolean): $CancellablePromise<models$0.InvitationTicket[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/governance/tickets${activeOnly ? '?active_only=true' : ''}`;
      const res = await httpReq<models$0.InvitationTicket[]>('GET', url);
      return res || [];
    })());
  },

  RevokeInvitationTicket(id: string): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/governance/tickets/${encodeURIComponent(id)}`));
  },

  ListUsersForGovernance(q: string, role: string, limit: number, offset: number): $CancellablePromise<models$0.UserProfile[]> {
    return this.ListUsers(q, role, limit, offset);
  },

  ListUsers(q: string, role: string, limit: number, offset: number): $CancellablePromise<models$0.UserProfile[]> {
    return toCancellable((async () => {
      const params = new URLSearchParams();
      if (q) params.set('q', q);
      if (role) params.set('role', role);
      if (limit > 0) params.set('limit', limit.toString());
      if (offset > 0) params.set('offset', offset.toString());
      const query = params.toString();
      const url = `${CORE_URL}/api/v1/governance/users${query ? `?${query}` : ''}`;
      const res = await httpReq<models$0.UserProfile[]>('GET', url);
      return res || [];
    })());
  },

  UpdateUserRole(id: string, role: models$0.AccountRole): $CancellablePromise<models$0.UserProfile | null> {
    return this.UpdateUserRoleOrStatus(id, { role });
  },

  UpdateUserStatus(id: string, isActive: boolean): $CancellablePromise<models$0.UserProfile | null> {
    return this.UpdateUserRoleOrStatus(id, { is_active: isActive });
  },

  UpdateUserRoleOrStatus(id: string, req: models$0.UpdateUserRoleOrStatusRequest): $CancellablePromise<models$0.UserProfile | null> {
    return toCancellable(
      httpReq<models$0.UserProfile>('PUT', `${CORE_URL}/api/v1/governance/users/${encodeURIComponent(id)}/access`, req)
    );
  },

  CreateCommunityReport(req: models$0.CreateReportRequest): $CancellablePromise<models$0.ContentReport | null> {
    return this.CreateReport(req);
  },

  CreateReport(req: models$0.CreateReportRequest): $CancellablePromise<models$0.ContentReport | null> {
    return toCancellable(httpReq<models$0.ContentReport>('POST', `${CORE_URL}/api/v1/reports`, req));
  },

  ListCommunityReports(status: string): $CancellablePromise<models$0.ContentReport[]> {
    return this.ListReports(status);
  },

  ListReports(status: string): $CancellablePromise<models$0.ContentReport[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/reports${status ? `?status=${encodeURIComponent(status)}` : ''}`;
      const res = await httpReq<models$0.ContentReport[]>('GET', url);
      return res || [];
    })());
  },

  ResolveCommunityReport(reportId: string, req: models$0.ResolveReportRequest): $CancellablePromise<models$0.ContentReport | null> {
    return this.ResolveReport(reportId, req);
  },

  ResolveReport(reportId: string, req: models$0.ResolveReportRequest): $CancellablePromise<models$0.ContentReport | null> {
    return toCancellable(
      httpReq<models$0.ContentReport>('PUT', `${CORE_URL}/api/v1/reports/${encodeURIComponent(reportId)}/resolve`, req)
    );
  },

  // ===========================================================================
  // MÓDULO 4: Diccionario, Regiones Dialectales y Conceptos Semánticos
  // ===========================================================================

  ListRegions(): $CancellablePromise<models$0.Region[]> {
    return toCancellable((async () => {
      const res = await httpReq<models$0.Region[]>('GET', `${CORE_URL}/api/v1/regions`);
      return res || [];
    })());
  },

  GetWordOfTheDay(dialect: string): $CancellablePromise<models$0.WordDetail | null> {
    const url = `${CORE_URL}/api/v1/words/word-of-the-day${dialect ? `?dialect=${encodeURIComponent(dialect)}` : ''}`;
    return toCancellable(httpReq<models$0.WordDetail>('GET', url));
  },

  SearchWords(q: string, dialect: string, limit: number): $CancellablePromise<models$0.WordSearchResult[]> {
    return toCancellable((async () => {
      const params = new URLSearchParams();
      params.set('q', q);
      if (dialect) params.set('dialect', dialect);
      if (limit > 0) params.set('limit', limit.toString());
      const res = await httpReq<models$0.WordSearchResult[]>('GET', `${CORE_URL}/api/v1/words/search?${params.toString()}`);
      return res || [];
    })());
  },

  ExploreWords(filter: models$0.WordExploreFilter): $CancellablePromise<models$0.PaginatedWords | null> {
    const params = new URLSearchParams();
    if (filter.q) params.set('q', filter.q);
    if (filter.region) params.set('region', filter.region);
    if (filter.pos) params.set('pos', filter.pos);
    if (filter.cefr) params.set('cefr', filter.cefr);
    if (filter.frequency) params.set('frequency', filter.frequency);
    params.set('include_taboo', String(filter.include_taboo ?? true));
    params.set('page', String(filter.page || 1));
    params.set('limit', String(filter.limit || 20));
    return toCancellable(httpReq<models$0.PaginatedWords>('GET', `${CORE_URL}/api/v1/words?${params.toString()}`));
  },

  GetWordDetail(id: number): $CancellablePromise<models$0.WordDetail | null> {
    return this.GetWordByID(id);
  },

  GetWordByID(id: number): $CancellablePromise<models$0.WordDetail | null> {
    return toCancellable(httpReq<models$0.WordDetail>('GET', `${CORE_URL}/api/v1/words/${id}`));
  },

  CreateWordProposal(req: models$0.CreateWordProposalRequest): $CancellablePromise<models$0.WordDetail | null> {
    return this.CreateWord(req);
  },

  CreateWord(req: models$0.CreateWordProposalRequest): $CancellablePromise<models$0.WordDetail | null> {
    return toCancellable(httpReq<models$0.WordDetail>('POST', `${CORE_URL}/api/v1/words`, req));
  },

  UpdateWordProposal(id: number, req: models$0.UpdateWordRequest): $CancellablePromise<models$0.WordDetail | null> {
    return this.UpdateWord(id, req);
  },

  UpdateWord(id: number, req: models$0.UpdateWordRequest): $CancellablePromise<models$0.WordDetail | null> {
    return toCancellable(httpReq<models$0.WordDetail>('PUT', `${CORE_URL}/api/v1/words/${id}`, req));
  },

  CreateSenseProposal(wordId: number, req: models$0.CreateSenseRequest): $CancellablePromise<models$0.SenseDetail | null> {
    return this.CreateSense(wordId, req);
  },

  CreateSense(wordId: number, req: models$0.CreateSenseRequest): $CancellablePromise<models$0.SenseDetail | null> {
    return toCancellable(httpReq<models$0.SenseDetail>('POST', `${CORE_URL}/api/v1/words/${wordId}/senses`, req));
  },

  UpdateSenseProposal(senseId: number, req: models$0.UpdateSenseRequest): $CancellablePromise<models$0.SenseDetail | null> {
    return this.UpdateSense(senseId, req);
  },

  UpdateSense(senseId: number, req: models$0.UpdateSenseRequest): $CancellablePromise<models$0.SenseDetail | null> {
    return toCancellable(httpReq<models$0.SenseDetail>('PUT', `${CORE_URL}/api/v1/senses/${senseId}`, req));
  },

  DeleteSense(senseId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/senses/${senseId}`));
  },

  AddSenseExample(senseId: number, req: models$0.CreateExampleInput): $CancellablePromise<models$0.Example | null> {
    return this.CreateExample(senseId, req);
  },

  CreateExample(senseId: number, req: models$0.CreateExampleInput): $CancellablePromise<models$0.Example | null> {
    return toCancellable(httpReq<models$0.Example>('POST', `${CORE_URL}/api/v1/senses/${senseId}/examples`, req));
  },

  UpdateExample(exampleId: number, req: models$0.CreateExampleInput): $CancellablePromise<models$0.Example | null> {
    return toCancellable(httpReq<models$0.Example>('PUT', `${CORE_URL}/api/v1/examples/${exampleId}`, req));
  },

  DeleteExample(exampleId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/examples/${exampleId}`));
  },

  ListSemanticConcepts(q: string): $CancellablePromise<models$0.SemanticConcept[]> {
    return this.ListConcepts(q);
  },

  ListConcepts(q: string): $CancellablePromise<models$0.SemanticConcept[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/concepts${q ? `?q=${encodeURIComponent(q)}` : ''}`;
      const res = await httpReq<models$0.SemanticConcept[]>('GET', url);
      return res || [];
    })());
  },

  CreateSemanticConcept(req: models$0.CreateConceptRequest): $CancellablePromise<models$0.SemanticConcept | null> {
    return this.CreateConcept(req);
  },

  CreateConcept(req: models$0.CreateConceptRequest): $CancellablePromise<models$0.SemanticConcept | null> {
    return toCancellable(httpReq<models$0.SemanticConcept>('POST', `${CORE_URL}/api/v1/concepts`, req));
  },

  LinkSenseToConcept(senseId: number, conceptId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('POST', `${CORE_URL}/api/v1/senses/${senseId}/concepts/${conceptId}`));
  },

  UnlinkSenseFromConcept(senseId: number, conceptId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/senses/${senseId}/concepts/${conceptId}`));
  },

  // ===========================================================================
  // MÓDULO 5: Curaduría Editorial por Quórum
  // ===========================================================================

  ListPendingQuorumQueue(region: string): $CancellablePromise<models$0.SenseQuorumItem[]> {
    return this.ListPendingQueue(region);
  },

  ListPendingQueue(region: string): $CancellablePromise<models$0.SenseQuorumItem[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/curation/queue${region ? `?region=${encodeURIComponent(region)}` : ''}`;
      const res = await httpReq<models$0.SenseQuorumItem[]>('GET', url);
      return res || [];
    })());
  },

  ListMyLexicalProposals(state: string): $CancellablePromise<models$0.SenseQuorumItem[]> {
    return this.ListMyProposals(state);
  },

  ListMyProposals(state: string): $CancellablePromise<models$0.SenseQuorumItem[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/curation/my-proposals${state ? `?state=${encodeURIComponent(state)}` : ''}`;
      const res = await httpReq<models$0.SenseQuorumItem[]>('GET', url);
      return res || [];
    })());
  },

  GetProposalDetail(senseId: number): $CancellablePromise<models$0.SenseQuorumDetail | null> {
    return toCancellable(httpReq<models$0.SenseQuorumDetail>('GET', `${CORE_URL}/api/v1/curation/senses/${senseId}`));
  },

  SubmitProposalReview(senseId: number, req: models$0.SubmitVoteRequest): $CancellablePromise<models$0.SenseQuorumDetail | null> {
    return this.SubmitVote(senseId, req);
  },

  SubmitVote(senseId: number, req: models$0.SubmitVoteRequest): $CancellablePromise<models$0.SenseQuorumDetail | null> {
    return toCancellable(
      httpReq<models$0.SenseQuorumDetail>('POST', `${CORE_URL}/api/v1/curation/senses/${senseId}/vote`, req)
    );
  },

  AdminForceResolveProposal(senseId: number, req: models$0.AdminOverrideStateRequest): $CancellablePromise<models$0.SenseDetail | null> {
    return this.AdminOverrideState(senseId, req);
  },

  AdminOverrideState(senseId: number, req: models$0.AdminOverrideStateRequest): $CancellablePromise<models$0.SenseDetail | null> {
    return toCancellable(
      httpReq<models$0.SenseDetail>('PUT', `${CORE_URL}/api/v1/curation/senses/${senseId}/state`, req)
    );
  },

  // ===========================================================================
  // MÓDULO 6: Studio Personal, SM-2 y Maestría
  // ===========================================================================

  ExplorePublicCollections(featured: boolean): $CancellablePromise<models$0.Collection[]> {
    return this.ListCollections(featured);
  },

  ListCollections(featured: boolean): $CancellablePromise<models$0.Collection[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/collections${featured ? '?featured=true' : ''}`;
      const res = await httpReq<models$0.Collection[]>('GET', url);
      return res || [];
    })());
  },

  GetCollectionDetail(id: number): $CancellablePromise<models$0.Collection | null> {
    return this.GetCollection(id);
  },

  GetCollection(id: number): $CancellablePromise<models$0.Collection | null> {
    return toCancellable(httpReq<models$0.Collection>('GET', `${CORE_URL}/api/v1/collections/${id}`));
  },

  ExplorePublicDecks(q: string, officialOnly: boolean): $CancellablePromise<models$0.Deck[]> {
    return this.ListDecks(q, officialOnly);
  },

  ListDecks(q: string, officialOnly: boolean): $CancellablePromise<models$0.Deck[]> {
    return toCancellable((async () => {
      const params = new URLSearchParams();
      if (q) params.set('q', q);
      if (officialOnly) params.set('official_only', 'true');
      const query = params.toString();
      const url = `${CORE_URL}/api/v1/decks${query ? `?${query}` : ''}`;
      const res = await httpReq<models$0.Deck[]>('GET', url);
      return res || [];
    })());
  },

  GetDeckDetail(id: number): $CancellablePromise<models$0.Deck | null> {
    return this.GetDeck(id);
  },

  GetDeck(id: number): $CancellablePromise<models$0.Deck | null> {
    return toCancellable(httpReq<models$0.Deck>('GET', `${CORE_URL}/api/v1/decks/${id}`));
  },

  ClonePublicDeck(deckId: number): $CancellablePromise<void> {
    return this.SaveDeck(deckId);
  },

  SaveDeck(deckId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('POST', `${CORE_URL}/api/v1/decks/${deckId}/save`));
  },

  UnsaveDeck(deckId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/decks/${deckId}/save`));
  },

  ListMyDecks(): $CancellablePromise<models$0.StudioDeckSummary[]> {
    return this.ListStudioDecks();
  },

  ListStudioDecks(): $CancellablePromise<models$0.StudioDeckSummary[]> {
    return toCancellable((async () => {
      const res = await httpReq<models$0.StudioDeckSummary[]>('GET', `${CORE_URL}/api/v1/studio/decks`);
      return res || [];
    })());
  },

  CreateDeck(req: models$0.CreateDeckRequest): $CancellablePromise<models$0.Deck | null> {
    return toCancellable(httpReq<models$0.Deck>('POST', `${CORE_URL}/api/v1/decks`, req));
  },

  UpdateDeck(id: number, req: models$0.UpdateDeckRequest): $CancellablePromise<models$0.Deck | null> {
    return toCancellable(httpReq<models$0.Deck>('PUT', `${CORE_URL}/api/v1/decks/${id}`, req));
  },

  DeleteDeck(id: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/decks/${id}`));
  },

  AddSenseToDeck(deckId: number, req: models$0.AddDeckItemRequest): $CancellablePromise<models$0.DeckItem | null> {
    return this.AddDeckItem(deckId, req);
  },

  AddDeckItem(deckId: number, req: models$0.AddDeckItemRequest): $CancellablePromise<models$0.DeckItem | null> {
    return toCancellable(httpReq<models$0.DeckItem>('POST', `${CORE_URL}/api/v1/decks/${deckId}/items`, req));
  },

  RemoveSenseFromDeck(deckId: number, senseId: number): $CancellablePromise<void> {
    return this.RemoveDeckItem(deckId, senseId);
  },

  RemoveDeckItem(deckId: number, senseId: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/decks/${deckId}/items/${senseId}`));
  },

  ListMyCollections(): $CancellablePromise<models$0.Collection[]> {
    return this.ListStudioCollections();
  },

  ListStudioCollections(): $CancellablePromise<models$0.Collection[]> {
    return toCancellable((async () => {
      const res = await httpReq<models$0.Collection[]>('GET', `${CORE_URL}/api/v1/studio/collections`);
      return res || [];
    })());
  },

  CreateCollection(req: models$0.CreateCollectionRequest): $CancellablePromise<models$0.Collection | null> {
    return toCancellable(httpReq<models$0.Collection>('POST', `${CORE_URL}/api/v1/collections`, req));
  },

  UpdateCollection(id: number, req: models$0.UpdateCollectionRequest): $CancellablePromise<models$0.Collection | null> {
    return toCancellable(httpReq<models$0.Collection>('PUT', `${CORE_URL}/api/v1/collections/${id}`, req));
  },

  DeleteCollection(id: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/collections/${id}`));
  },

  SaveCollection(id: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('POST', `${CORE_URL}/api/v1/collections/${id}/save`));
  },

  UnsaveCollection(id: number): $CancellablePromise<void> {
    return toCancellable(httpReq<void>('DELETE', `${CORE_URL}/api/v1/collections/${id}/save`));
  },

  AddDeckToCollection(collectionID: number, deckID: number): $CancellablePromise<models$0.Collection | null> {
    return toCancellable((async () => {
      const col = await this.GetCollection(collectionID);
      if (!col) throw new Error('Colección no encontrada');
      const existing = (col.decks || []).map((d) => d.id_deck);
      if (!existing.includes(deckID)) {
        existing.push(deckID);
      }
      return await this.UpdateCollection(collectionID, {
        title: col.title,
        description: col.description,
        is_public: col.is_public,
        is_featured: col.is_featured,
        is_official: col.is_official,
        deck_ids: existing,
      });
    })());
  },

  RemoveDeckFromCollection(collectionID: number, deckID: number): $CancellablePromise<models$0.Collection | null> {
    return toCancellable((async () => {
      const col = await this.GetCollection(collectionID);
      if (!col) throw new Error('Colección no encontrada');
      const filtered = (col.decks || []).map((d) => d.id_deck).filter((id) => id !== deckID);
      return await this.UpdateCollection(collectionID, {
        title: col.title,
        description: col.description,
        is_public: col.is_public,
        is_featured: col.is_featured,
        is_official: col.is_official,
        deck_ids: filtered,
      });
    })());
  },

  GetStudySessionQueue(deckId: number, limit: number): $CancellablePromise<models$0.StudySessionPayload | null> {
    return this.GetStudySession(deckId, limit);
  },

  GetStudySession(deckId: number, limit: number): $CancellablePromise<models$0.StudySessionPayload | null> {
    const params = new URLSearchParams();
    if (deckId > 0) params.set('deck_id', deckId.toString());
    if (limit > 0) params.set('limit', limit.toString());
    const query = params.toString();
    const url = `${CORE_URL}/api/v1/studio/study-session${query ? `?${query}` : ''}`;
    return toCancellable(httpReq<models$0.StudySessionPayload>('GET', url));
  },

  SubmitCardReview(req: models$0.SubmitReviewRequest): $CancellablePromise<models$0.CardReviewResult | null> {
    return this.SubmitReview(req);
  },

  SubmitReview(req: models$0.SubmitReviewRequest): $CancellablePromise<models$0.CardReviewResult | null> {
    return toCancellable(httpReq<models$0.CardReviewResult>('POST', `${CORE_URL}/api/v1/studio/study-session/review`, req));
  },

  SubmitStudyReview(req: models$0.SubmitReviewRequest): $CancellablePromise<models$0.CardReviewResult | null> {
    return this.SubmitReview(req);
  },

  GetMasterySummary(): $CancellablePromise<models$0.MasteryDashboard | null> {
    return this.GetMasteryDashboard();
  },

  GetMasteryDashboard(): $CancellablePromise<models$0.MasteryDashboard | null> {
    return toCancellable(httpReq<models$0.MasteryDashboard>('GET', `${CORE_URL}/api/v1/studio/mastery`));
  },

  ListMasteryItems(stage: string): $CancellablePromise<models$0.UserMastery[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/studio/mastery/items${stage ? `?stage=${encodeURIComponent(stage)}` : ''}`;
      const res = await httpReq<models$0.UserMastery[]>('GET', url);
      return res || [];
    })());
  },

  GetMyXPTransactions(limit: number): $CancellablePromise<models$0.XPTransactionItem[]> {
    return this.ListXPLedger(limit);
  },

  ListXPLedger(limit: number): $CancellablePromise<models$0.XPTransactionItem[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/studio/xp-ledger${limit > 0 ? `?limit=${limit}` : ''}`;
      const res = await httpReq<models$0.XPTransactionItem[]>('GET', url);
      return res || [];
    })());
  },

  GetMyAchievements(): $CancellablePromise<models$0.AchievementStatus[]> {
    return toCancellable((async () => {
      const dash = await this.GetMasteryDashboard();
      return dash?.achievements || [];
    })());
  },

  // ===========================================================================
  // MÓDULO 7: Arena Multijugador y WebSockets
  // ===========================================================================

  CreateArenaMatch(req: models$0.CreateMatchRequest): $CancellablePromise<models$0.MatchPackage | null> {
    return this.PrepareMatch(req);
  },

  PrepareMatch(req: models$0.CreateMatchRequest): $CancellablePromise<models$0.MatchPackage | null> {
    return toCancellable(httpReq<models$0.MatchPackage>('POST', `${CORE_URL}/api/v1/arena/matches/prepare`, req));
  },

  JoinArenaMatch(roomCode: string): $CancellablePromise<models$0.MatchPackage | null> {
    return toCancellable(Promise.resolve({ room_code: roomCode } as models$0.MatchPackage));
  },

  GetMyMatchHistory(limit: number): $CancellablePromise<models$0.MatchSummary[]> {
    return this.ListMatchHistory(limit);
  },

  ListMatchHistory(limit: number): $CancellablePromise<models$0.MatchSummary[]> {
    return toCancellable((async () => {
      const url = `${CORE_URL}/api/v1/arena/matches/history${limit > 0 ? `?limit=${limit}` : ''}`;
      const res = await httpReq<models$0.MatchSummary[]>('GET', url);
      return res || [];
    })());
  },

  GetMatchReplayDetail(matchId: string): $CancellablePromise<models$0.MatchReplayDetail | null> {
    return this.GetMatchReplay(matchId);
  },

  GetMatchReplay(matchId: string): $CancellablePromise<models$0.MatchReplayDetail | null> {
    return toCancellable(httpReq<models$0.MatchReplayDetail>('GET', `${CORE_URL}/api/v1/arena/matches/${encodeURIComponent(matchId)}`));
  },

  ConnectArenaRoom(roomCode: string): $CancellablePromise<void> {
    return toCancellable((async () => {
      if (activeWebSocket) {
        activeWebSocket.close();
        activeWebSocket = null;
      }

      const token = await getSessionToken();
      if (!token) {
        throw new Error('se requiere una sesión activa para ingresar a la arena');
      }

      let wsBase = REALTIME_URL;
      if (wsBase.startsWith('http://')) {
        wsBase = 'ws://' + wsBase.slice(7);
      } else if (wsBase.startsWith('https://')) {
        wsBase = 'wss://' + wsBase.slice(8);
      }

      const wsURL = `${wsBase}/ws/arena/${encodeURIComponent(roomCode)}?token=${encodeURIComponent(token)}`;
      const ws = new WebSocket(wsURL);
      activeWebSocket = ws;

      await new Promise<void>((resolve, reject) => {
        ws.onopen = () => resolve();
        ws.onerror = (err) => reject(new Error('error al conectar con la sala de arena'));
      });

      ws.onmessage = (event) => {
        try {
          const evt = JSON.parse(event.data) as models$0.WSEvent;
          ArenaEventBus.dispatchBrowserMessage(evt);
        } catch (e) {
          console.error('Error deserializando mensaje WS en ArenaEventBus:', e);
        }
      };

      ws.onclose = () => {
        if (activeWebSocket === ws) {
          activeWebSocket = null;
        }
      };
    })());
  },

  ConnectArena(roomCode: string): $CancellablePromise<void> {
    return this.ConnectArenaRoom(roomCode);
  },

  StartArenaMatch(): $CancellablePromise<void> {
    return toCancellable((async () => {
      if (!activeWebSocket || activeWebSocket.readyState !== WebSocket.OPEN) {
        throw new Error('no hay conexión activa con la sala de arena');
      }
      activeWebSocket.send(JSON.stringify({ type: 'START_MATCH' }));
    })());
  },

  SubmitArenaAnswer(roundNumber: number, selectedOption: string): $CancellablePromise<void> {
    return toCancellable((async () => {
      if (!activeWebSocket || activeWebSocket.readyState !== WebSocket.OPEN) {
        throw new Error('no hay conexión activa con la sala de arena');
      }
      activeWebSocket.send(
        JSON.stringify({
          type: 'SUBMIT_ANSWER',
          payload: {
            round_number: roundNumber,
            selected_option: selectedOption,
          },
        })
      );
    })());
  },

  LeaveArenaRoom(): $CancellablePromise<void> {
    return toCancellable((async () => {
      if (activeWebSocket) {
        activeWebSocket.close();
        activeWebSocket = null;
      }
    })());
  },

  LeaveArena(): $CancellablePromise<void> {
    return this.LeaveArenaRoom();
  },

  // ===========================================================================
  // MÓDULO 8: Multimedia Seguro
  // ===========================================================================

  UploadImage(category: string, filename: string, fileBytes: string): $CancellablePromise<models$0.ImageUploadResponse | null> {
    return toCancellable((async () => {
      const formData = new FormData();
      formData.append('category', category);

      let blob: Blob;
      try {
        const cleanBase64 = fileBytes.replace(/^data:image\/\w+;base64,/, '');
        const byteChars = atob(cleanBase64);
        const byteNumbers = new Array(byteChars.length);
        for (let i = 0; i < byteChars.length; i++) {
          byteNumbers[i] = byteChars.charCodeAt(i);
        }
        blob = new Blob([new Uint8Array(byteNumbers)]);
      } catch {
        blob = new Blob([fileBytes]);
      }

      formData.append('file', blob, filename || 'image.webp');
      return await httpReq<models$0.ImageUploadResponse>('POST', `${IMAGE_URL}/api/v1/images/upload`, formData);
    })());
  },
};
