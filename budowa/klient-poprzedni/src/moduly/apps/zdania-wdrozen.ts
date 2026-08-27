import { AppDeployStatus, type AppDeployment } from '../../../../shared/contract';
import type { RachunekRamek } from './zbior-budowy';

/**
 * Zamówienie wdrożenia wysłane z okna Deployment Panelu, zestawiane potem
 * z odpowiedzią rdzenia. Plik składa zdania panelu z liczb i pól, bez kanału
 * i bez elementu, więc czyta się go w całości bez czytania okna.
 */
export interface Zamowienie {
  srodowisko: string;
  strategia: string;
  /** Wdrożenie, do którego cofamy; pusty łańcuch znaczy „wdrożenie w przód". */
  cofnijDo: string;
}

/**
 * Wylicza, czym przebieg oddany przez rdzeń różni się od zamówienia; zgodność
 * daje pusty łańcuch. Cofnięcie potwierdza osobne pole `rolledBackFromId`
 * odpowiedzi, a nie samo zamówienie wysłane z okna.
 */
export function rozbieznoscZlecenia(zamowienie: Zamowienie, przebieg: AppDeployment): string {
  const rozejscia: string[] = [];
  if (przebieg.environment !== zamowienie.srodowisko) {
    rozejscia.push(
      `zamówiono środowisko ${zamowienie.srodowisko}, rdzeń oddał ${przebieg.environment}`,
    );
  }
  if (przebieg.strategy !== zamowienie.strategia) {
    rozejscia.push(`zamówiono strategię ${zamowienie.strategia}, rdzeń oddał ${przebieg.strategy}`);
  }
  const cofnietoDo = przebieg.rolledBackFromId ?? '';
  if (cofnietoDo !== zamowienie.cofnijDo) {
    rozejscia.push(
      zamowienie.cofnijDo === ''
        ? `nie zamawiano cofnięcia, a rdzeń oddał przebieg cofnięty z ${cofnietoDo}`
        : `zamówiono cofnięcie do ${zamowienie.cofnijDo}, a rdzeń oddał ` +
          (cofnietoDo === '' ? 'przebieg bez wskazania wdrożenia źródłowego' : cofnietoDo),
    );
  }
  return rozejscia.join('; ');
}

/**
 * Składa zdanie o powodzie pustego wykazu wdrożeń. Pustkę po udanym odczycie
 * komendą `apps.deployment.list` oddziela od chwili przed odczytem, a tę
 * rozstrzyga rachunkiem ramek przyjętych przez okno.
 */
export function zdaniePustkiWdrozen(ramki: RachunekRamek, czytane: boolean): string {
  if (czytane) {
    return (
      'Rdzeń odpowiedział na odczyt wdrożeń i nie zna ani jednego wdrożenia tego okna. ' +
      'To jest pustka POTWIERDZONA, nie brak odczytu — wykaz wypełni się, gdy pierwsze ' +
      'wdrożenie ruszy albo gdy przyjdzie ramka apps.build.changed.'
    );
  }
  if (ramki.wszystkie === 0) {
    return (
      'Brak wdrożeń, a wykaz nie był jeszcze czytany z rdzenia — naciśnij „Odczytaj wdrożenia". ' +
      'Zdarzenie apps.build.changed nie przyszło ani razu i żadne zlecenie z tego okna nie ' +
      'wróciło jeszcze z odpowiedzią, więc pustka znaczy tu „nic jeszcze nie przyszło", ' +
      'a nie „rdzeń nic nie ma".'
    );
  }
  const ile = `Brak wdrożeń mimo ${ramki.wszystkie} ramek apps.build.changed`;
  if (ramki.zWdrozeniem === 0) {
    return (
      `${ile}: w żadnej z nich nie było pola wdrożenia, a wykazu nie czytano jeszcze ` +
      'z rdzenia — naciśnij „Odczytaj wdrożenia".'
    );
  }
  return (
    `${ile}, z których ${ramki.zWdrozeniem} niosło wdrożenie. Wykaz mimo to jest pusty, ` +
    'choć nic go nie opróżnia — to sprzeczność wewnątrz okna, nie stan rdzenia; zgłoś ją.'
  );
}

/**
 * Rozstrzyga, czy przebieg wdrożenia się domknął. Stany końcowe bierze
 * z wyliczenia kontraktu `AppDeployStatus`, a nie z literałów, więc nowy stan
 * przerwie kompilację, zamiast po cichu wypaść z rozpoznania.
 */
export function czyStanKoncowy(stan: string): boolean {
  return (
    stan === AppDeployStatus.Succeeded ||
    stan === AppDeployStatus.Failed ||
    stan === AppDeployStatus.RolledBack
  );
}

/**
 * Składa zdanie o niespełnionym warunku prac w warsztatach. Moduł wie o nich
 * tyle, ile potwierdził rdzeń zapisem albo odczytem architektury, a brak
 * potwierdzenia daje ostrzeżenie, nie blokadę wdrożenia.
 */
export function warunekWarsztatow(maArchitekture: boolean): string {
  if (maArchitekture) return '';
  return (
    ' Uwaga: rdzeń nie potwierdził jeszcze żadnego zapisu architektury ani warsztatu, ' +
    'więc warunek „aktywny po zakończeniu prac w warsztatach” nie jest spełniony.'
  );
}
