import type { AuthSession } from '../../../shared/contract';

/** Klucz zapisu sesji bramki, pod którym moduł przechowuje ją w pamięci trwałej albo w pamięci okna przeglądarki zależnie od wyboru trwałości logowania. */
export const KLUCZ_SESJI = 'danaco-console.sesja-bramki';

/** Zwraca pamięć trwałą przeglądarki, w której zapis sesji przeżywa zamknięcie aplikacji, albo pustą wartość, gdy dostęp jest niemożliwy. */
function pamiecTrwala(): Storage | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

/** Zwraca pamięć okna przeglądarki, w której zapis sesji ginie wraz z jego zamknięciem, albo pustą wartość, gdy dostęp jest niemożliwy. */
function pamiecOkna(): Storage | null {
  try {
    return window.sessionStorage;
  } catch {
    return null;
  }
}

/** Sprawdza, czy odczytana wartość ma kształt sesji zgodny z kontraktem, zawierający niepusty token oraz liczbowy czas wygaśnięcia. */
function czySesja(wartosc: unknown): wartosc is AuthSession {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const zapis = wartosc as Record<string, unknown>;
  return (
    typeof zapis['token'] === 'string' &&
    zapis['token'] !== '' &&
    typeof zapis['expiresAt'] === 'number'
  );
}

/** Odczytuje sesję z jednego wskazanego magazynu; zapis nieczytelny albo pusty jest traktowany tak samo jak jego brak. */
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

/** Odczytuje sesję zapisaną przy poprzednim wejściu, sprawdzając kolejno pamięć trwałą i pamięć okna; brak albo zapis nieczytelny daje pustą wartość. */
export function odczytajSesje(): AuthSession | null {
  return odczytajZ(pamiecTrwala()) ?? odczytajZ(pamiecOkna());
}

/** Zapisuje sesję po wejściu albo po przedłużeniu do jednego magazynu wskazanego parametrem trwałości, czyszcząc przy tym zapis z magazynu drugiego. */
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

/** Kasuje zapis sesji z obu magazynów przeglądarki jednocześnie, tak aby sesja martwa nie zalegała w żadnym z nich. */
export function skasujSesje(): void {
  try {
    pamiecTrwala()?.removeItem(KLUCZ_SESJI);
    pamiecOkna()?.removeItem(KLUCZ_SESJI);
  } catch {
    // Zapis nieusuwalny nie szkodzi: następne przedłużenie i tak odmówi.
  }
}

/** Sprawdza, czy zapisana sesja leży w pamięci trwałej, co odpowiada zaznaczonemu polu trwałości logowania na ekranie wejścia. */
export function sesjaTrwala(): boolean {
  return odczytajZ(pamiecTrwala()) !== null;
}

/** Buduje zdanie o ważności sesji na podstawie czasu wygaśnięcia zwróconego przez rdzeń, pokazujące godzinę w dniu bieżącym albo pełną datę w dniach kolejnych. */
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
