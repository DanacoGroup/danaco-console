import {
  Command,
  ConfigScope,
  type IsolationPolicyPreviewRequest,
  type IsolationPolicyPreviewResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

/**
 * Polityka izolacji obowiązująca — jedyna droga kontraktu do „statusu sandboxu
 * wykonania" i „podglądu macierzy izolacji" w Permissions & Trust Center.
 *
 * Granica jest tu istotna i okno musi ją nazwać: `isolation.policy.preview`
 * oddaje politykę PLATFORMY rozstrzygniętą po ośmiu poziomach zasięgu, a nie
 * izolację nadaną pojedynczemu rozszerzeniu. Osiem zakresów technicznych
 * (katalog roboczy, środowisko procesu, dostęp sieciowy, odczyt i zapis plików,
 * konto i token, model procesu, serwer wykonania, katalog danych modelu) opisuje
 * warunki, w jakich kod się wykonuje — i w tych samych warunkach wykona się kod
 * rozszerzenia. Nie jest to jednak deklaracja uprawnień rozszerzenia, bo takiej
 * kontrakt nie niesie; okno mówi o tym wprost zamiast podstawiać jedno za drugie.
 *
 * Poziom zasięgu bierzemy najwęższy, jaki moduł zna: podanie okna modułu każe
 * rdzeniowi rozstrzygnąć dziedziczenie aż do niego, więc polityka jest tą,
 * która naprawdę obowiązuje pracy prowadzonej w tym oknie. Bez okna pytamy
 * o poziom modułu — wtedy odpowiedź opisuje warunki wspólne wszystkim jego
 * oknom, co jest zdaniem prawdziwym, tyle że szerszym.
 *
 * Warstwy nie podajemy: pominięta znaczy „warstwa obowiązująca", a wskazanie
 * którejkolwiek byłoby cudzą decyzją przebraną za odczyt.
 */
export interface ZrodloIzolacjiApps {
  /** `isolation.policy.preview` — polityka obowiązująca oknu albo modułowi. */
  polityka(idOkna: string): Promise<Wynik<IsolationPolicyPreviewResponse>>;
}

/** Kod modułu jako byt poziomu `module`; ten sam, którym rdzeń zna moduł. */
export function utworzZrodloIzolacjiApps(kanal: Kanal, kodModulu: string): ZrodloIzolacjiApps {
  const wywolaj: WywolanieApps = utworzWywolanieApps(kanal);

  return {
    async polityka(idOkna) {
      const zadanie: IsolationPolicyPreviewRequest =
        idOkna === ''
          ? { scope: ConfigScope.Module, scopeId: kodModulu }
          : { scope: ConfigScope.Window, scopeId: idOkna, windowId: idOkna };
      return sprawdzKsztalt(
        await wywolaj(Command.IsolationPolicyPreview, zadanie),
        Command.IsolationPolicyPreview,
        (tresc) => czyObiekt(tresc.policy),
      );
    },
  };
}
