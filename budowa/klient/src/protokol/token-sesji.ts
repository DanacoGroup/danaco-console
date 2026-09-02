// Token sesji bramki stoi w warstwie protokołu, bo powitanie ponawiane po
// zerwaniu musi go nieść z powrotem.
import { Command, EnvelopeStatus, type AuthSession } from '../../../shared/contract.ts';
import type { Kanal } from './kanal.ts';

let token = '';
// `auth.changed` o unieważnieniu wskazuje urządzenie, nie token.
let urzadzenie = '';

// Rdzeń wydaje token także przy przedłużeniu sesji, o którym okno nie wie;
// odpowiedź powtarza typ żądania, więc rozpoznaje się ją po nazwie komendy.
const KOMENDY_SESJI: readonly string[] = [
  Command.AuthLogin,
  Command.AuthVerify,
  Command.AuthTokenRefresh,
];

/** Token bieżącej sesji bramki; pustka znaczy, że Operator nie przeszedł wejścia. */
export function tokenSesji(): string {
  return token;
}

/** Urządzenie sesji bieżącej według rdzenia; pustka, gdy rdzeń go nie nazwał. */
export function urzadzenieSesji(): string {
  return urzadzenie;
}

/** Zapomina token; rdzeń unieważnia go przy zmianie hasła i przy odzyskaniu konta. */
export function zapomnijTokenSesji(): void {
  token = '';
  urzadzenie = '';
}

/** Zapamiętuje token z każdej odpowiedzi rdzenia, która sesję bramki wydaje albo przedłuża. */
export function pilnujTokenu(kanal: Kanal): void {
  kanal.naDowolny((koperta) => {
    if (koperta.status !== EnvelopeStatus.Ok) return;
    if (!KOMENDY_SESJI.includes(koperta.type)) return;
    const wydana = sesjaZTresci(koperta.payload);
    if (wydana === null) return;
    token = wydana.token;
    urzadzenie = wydana.deviceId ?? '';
  });
}

/** Sesja bramki z treści odpowiedzi; pustka, gdy odpowiedź jej nie niesie. */
function sesjaZTresci(tresc: unknown): AuthSession | null {
  if (typeof tresc !== 'object' || tresc === null) return null;
  const sesja = (tresc as { session?: unknown }).session;
  if (typeof sesja !== 'object' || sesja === null) return null;
  const wydany = (sesja as { token?: unknown }).token;
  if (typeof wydany !== 'string' || wydany === '') return null;
  return sesja as AuthSession;
}
