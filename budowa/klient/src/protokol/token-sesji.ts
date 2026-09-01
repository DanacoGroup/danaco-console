import { Command, EnvelopeStatus, type AuthSession } from '../../../shared/contract.ts';
import type { Kanal } from './kanal.ts';

/**
 * Token sesji bramki wydany przez rdzeń.
 *
 * Token stoi w warstwie protokołu, nie w wiązaniu okna wejścia: rdzeń wiąże
 * sesję bramki z gniazdem przy powitaniu, więc powitanie ponawiane po zerwaniu
 * połączenia musi go nieść z powrotem. Zamknięty w zasięgu jednego wiązania
 * byłby niedostępny warstwie, która to powitanie ponawia.
 */
let token = '';
/* Urządzenie, na którym sesja trwa; `auth.changed` o unieważnieniu sesji
   wskazuje urządzenie, więc po nim klient poznaje, czy to jego token. */
let urzadzenie = '';

/*
Komendy, których odpowiedź niesie sesję bramki. Token bierze się z odpowiedzi
rdzenia, a nie z ręki wiązania: rdzeń wydaje go także przy przedłużeniu sesji,
o którym okno wejścia nie wie, a powitanie po zerwaniu ma nieść ostatni wydany.
Odpowiedź powtarza typ żądania, więc rozpoznaje się ją po nazwie komendy.
*/
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
