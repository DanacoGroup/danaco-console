import type { AuthSession } from '../../../shared/contract';

/**
 * Trwały zapis sesji bramki w pamięci przeglądarki, wzorem `motyw/motyw.ts`.
 *
 * Transportem jest jedno gniazdo WebSocket, nie seria żądań HTTP, więc
 * ciasteczka nie ma: token wraca z `auth.login` i tylko klient może go
 * przechować między uruchomieniami. Bez zapisu każdy start aplikacji wymagałby
 * hasła, a bramka ma stać przy wejściu raz, nie przy każdym otwarciu okna.
 *
 * W zapisie leży sesja w kształcie kontraktu (`AuthSession`): token, czas
 * wygaśnięcia, metoda. Token jest poświadczeniem na okaziciela — rdzeń trzyma
 * w bazie wyłącznie jego skrót, więc jedyną kopią tokenu jest właśnie ten
 * zapis. `localStorage` interfejsu podawanego z `127.0.0.1` i z powłoki
 * natywnej jest magazynem lokalnym maszyny Operatora, czyli tym samym progiem
 * zaufania, na którym stoi plik sejfu rdzenia obok bazy.
 *
 * Magazyny są dwa, bo Operator ma wybór. Pole „Nie wyloguj mnie" na ekranie
 * logowania rozstrzyga, gdzie sesja wyląduje: zaznaczone kładzie ją
 * w `localStorage` i sesja przeżywa zamknięcie aplikacji, niezaznaczone —
 * w `sessionStorage`, skąd ginie razem z oknem.
 *
 * Skutek jest dwustronny: ta sama wartość idzie żądaniem `auth.login` jako
 * `keepSignedIn` i rozstrzyga w rdzeniu o trwaniu sesji — 365 dni zamiast doby
 * roboczej (`adapter_modul_auth.go`). Trwanie jest zapisane przy wierszu sesji,
 * więc `auth.token.refresh` go nie ścina. Zapis w przeglądarce bez tego byłby
 * obietnicą, której rdzeń nie dotrzymuje.
 *
 * Odczyt pyta obu magazynów, w kolejności od trwałego. Inaczej zmiana
 * rozstrzygnięcia między jednym a drugim uruchomieniem zostawiałaby sesję
 * niewidoczną, a wyglądałoby to jak jej wygaśnięcie. Kasowanie czyści oba
 * z tego samego powodu.
 *
 * Żaden błąd pamięci nie zatrzymuje uruchomienia: brak magazynu, zapis
 * nieczytelny i kształt spoza kontraktu znaczą to samo — sesji zapisanej nie ma
 * i Operator wchodzi hasłem. O tym, czy sesja jeszcze żyje, rozstrzyga wyłącznie
 * rdzeń odpowiedzią na `auth.token.refresh`; ten plik nie porównuje czasów
 * i niczego nie unieważnia sam.
 */

/** Klucz zapisu sesji bramki w pamięci trwałej przeglądarki. */
export const KLUCZ_SESJI = 'danaco-console.sesja-bramki';

/** Pamięć trwała — przeżywa zamknięcie aplikacji. */
function pamiecTrwala(): Storage | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

/** Pamięć okna — ginie razem z nim. */
function pamiecOkna(): Storage | null {
  try {
    return window.sessionStorage;
  } catch {
    return null;
  }
}

/** Czy odczytana wartość ma kształt sesji kontraktu. */
function czySesja(wartosc: unknown): wartosc is AuthSession {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const zapis = wartosc as Record<string, unknown>;
  return (
    typeof zapis['token'] === 'string' &&
    zapis['token'] !== '' &&
    typeof zapis['expiresAt'] === 'number'
  );
}

/** Odczyt z jednego magazynu; zapis nieczytelny znaczy to samo co brak. */
function odczytajZ(magazyn: Storage | null): AuthSession | null {
  try {
    const zapis = magazyn?.getItem(KLUCZ_SESJI);
    if (zapis === null || zapis === undefined || zapis === '') return null;
    const wartosc: unknown = JSON.parse(zapis);
    return czySesja(wartosc) ? wartosc : null;
  } catch {
    return null;
  }
}

/** Odczytuje sesję zapisaną przy poprzednim wejściu. Brak albo zapis nieczytelny → null. */
export function odczytajSesje(): AuthSession | null {
  return odczytajZ(pamiecTrwala()) ?? odczytajZ(pamiecOkna());
}

/**
 * Zapisuje sesję po wejściu albo po przedłużeniu.
 *
 * @param niewylogowuj czy sesja ma przeżyć zamknięcie aplikacji — rozstrzyga
 *   pole „Nie wyloguj mnie". Zapis idzie do jednego magazynu, a drugi jest
 *   czyszczony: sesja w dwóch miejscach naraz znaczyłaby, że odznaczenie pola
 *   niczego nie cofa.
 */
export function zapiszSesje(sesja: AuthSession, niewylogowuj: boolean): void {
  const cel = niewylogowuj ? pamiecTrwala() : pamiecOkna();
  const drugi = niewylogowuj ? pamiecOkna() : pamiecTrwala();
  try {
    drugi?.removeItem(KLUCZ_SESJI);
    cel?.setItem(KLUCZ_SESJI, JSON.stringify(sesja));
  } catch {
    // Brak zapisu nie odbiera wejścia — następne uruchomienie poprosi o hasło.
  }
}

/** Kasuje zapis z obu magazynów — sesja martwa nie zalega w żadnym. */
export function skasujSesje(): void {
  try {
    pamiecTrwala()?.removeItem(KLUCZ_SESJI);
    pamiecOkna()?.removeItem(KLUCZ_SESJI);
  } catch {
    // Zapis nieusuwalny nie szkodzi: następne przedłużenie i tak odmówi.
  }
}

/**
 * Czy zapisana sesja leży w pamięci trwałej — stan wyjściowy pola
 * „Nie wyloguj mnie".
 *
 * Jedna prawda o nastawie: pole nie pamięta własnego zaznaczenia osobnym
 * zapisem, tylko czyta miejsce, w którym sesja naprawdę leży. Osobny zapis
 * rozjechałby się z magazynem przy pierwszym czyszczeniu pamięci i pokazywałby
 * „nie wyloguj mnie" nad sesją, która ginie z oknem.
 */
export function sesjaTrwala(): boolean {
  return odczytajZ(pamiecTrwala()) !== null;
}

/**
 * Zdanie o ważności sesji pokazywane Operatorowi przed rozpoczęciem pracy.
 *
 * Czas idzie z odpowiedzi rdzenia (`expiresAt`), nie ze stałej klienta —
 * długość życia sesji zna wyłącznie rdzeń. Data pełna pojawia się
 * tylko wtedy, gdy wygaśnięcie wypada innego dnia; w dniu bieżącym wystarcza
 * godzina.
 */
export function opisWaznosci(expiresAt: number, teraz: number = Date.now()): string {
  const wygasa = new Date(expiresAt);
  const dzis = new Date(teraz);
  const tenSamDzien =
    wygasa.getFullYear() === dzis.getFullYear() &&
    wygasa.getMonth() === dzis.getMonth() &&
    wygasa.getDate() === dzis.getDate();
  const godzina = wygasa.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
  if (tenSamDzien) return `sesja ważna do ${godzina}`;
  const data = wygasa.toLocaleDateString('pl-PL', { day: 'numeric', month: 'long' });
  return `sesja ważna do ${data}, ${godzina}`;
}
