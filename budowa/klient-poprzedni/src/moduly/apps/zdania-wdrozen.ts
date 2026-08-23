import { AppDeployStatus, type AppDeployment } from '../../../../shared/contract';
import type { RachunekRamek } from './zbior-budowy';

/**
 * Zdania Deployment Panelu o wdrożeniu — rozbieżność zamówienia z odpowiedzią
 * i powód pustego wykazu.
 *
 * Stoi osobno od okna. Okno prowadzi rozmowę z rdzeniem: zbiera pola, wysyła,
 * przyjmuje odpowiedź, nazywa stan. Składanie zdań z liczb i pól to inna
 * czynność — czysta, bez kanału i bez elementu — i dzięki rozdziałowi da się ją
 * przeczytać w całości bez czytania okna.
 */

/** Zamówienie wysłane z okna — do zestawienia z odpowiedzią rdzenia. */
export interface Zamowienie {
  srodowisko: string;
  strategia: string;
  /** Wdrożenie, do którego cofamy; pusty łańcuch znaczy „wdrożenie w przód". */
  cofnijDo: string;
}

/**
 * Czym odpowiedź rdzenia różni się od zamówienia — pusty łańcuch, gdy niczym.
 *
 * Zdanie potwierdzające mówi o tym, co rdzeń oddał, a nie o tym, co okno
 * wysłało. Samo wypisanie pól odpowiedzi to za mało: zamówiono konkretne
 * środowisko i konkretną czynność, więc gdy rdzeń odda co innego, rozejście
 * pada wprost, a nie w drobnym druku. Cofnięcie sprawdzamy osobno, bo słowo
 * „Cofnięcie wdrożenia" w zdaniu bierze się z zamówienia — potwierdza je
 * dopiero pole `rolledBackFromId` w odpowiedzi.
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
 * Zdanie o powodzie pustego wykazu wdrożeń.
 *
 * Odczyt rozstrzyga pierwszy: pustka po udanym `apps.deployment.list` znaczy
 * rzecz ostateczną — rdzeń przejrzał bazę i wdrożeń tego okna nie ma. To zdanie
 * innej wagi niż „jeszcze nic nie przyszło" i te dwa stany nie są zlewane.
 *
 * Pozostałe trzy gałęzie dotyczą chwili, gdy nikt jeszcze nie pytał. Zero ramek
 * znaczy „nic nie przyszło i nie pytano". Ramki bez pola wdrożenia znaczą
 * „przyszło, ale wdrożeń w tym nie było". Ramki z wdrożeniem przy pustym
 * wykazie są sprzecznością wewnątrz okna, bo wykaz niczego nie zdejmuje; zdanie
 * nazywa ją wprost, zamiast zaokrąglić do jednej z pozostałych.
 *
 * Żadna gałąź nie orzeka o tym, czego rdzeń nie robi — zdanie mówi wyłącznie
 * o tym, co padło albo nie padło w tym oknie.
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
 * Czy przebieg wdrożenia się domknął.
 *
 * Wyliczenie stanów bierzemy z kontraktu, a nie z literałów — nowy stan
 * dopisany do `AppDeployStatus` przerwie kompilację tutaj, zamiast po cichu
 * wypaść z rozpoznania.
 */
export function czyStanKoncowy(stan: string): boolean {
  return (
    stan === AppDeployStatus.Succeeded ||
    stan === AppDeployStatus.Failed ||
    stan === AppDeployStatus.RolledBack
  );
}

/**
 * Zdanie o brakującym warunku merytorycznym — prace w warsztatach.
 *
 * Warunek sprawdzamy, a nie zakładamy: moduł wie o pracach w warsztatach tyle,
 * ile potwierdził rdzeń — czyli tyle, ile widać w architekturze, którą oddał
 * zapis albo odczyt. Brak potwierdzenia daje ostrzeżenie, nie blokadę.
 */
export function warunekWarsztatow(maArchitekture: boolean): string {
  if (maArchitekture) return '';
  return (
    ' Uwaga: rdzeń nie potwierdził jeszcze żadnego zapisu architektury ani warsztatu, ' +
    'więc warunek „aktywny po zakończeniu prac w warsztatach” nie jest spełniony.'
  );
}
