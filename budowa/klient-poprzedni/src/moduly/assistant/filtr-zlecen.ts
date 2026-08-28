import { AssistantActionStatus, type AssistantAction } from '../../../../shared/contract';
import { wybor } from '../../modele/kontrolki-formularza';
import { NAZWY_FILTRA, POZYCJE_FILTRA, type KodFiltraZlecen } from './etykiety-assistant';

/**
 * Filtr statusu wykazu zleceń Actions Monitora. Zawężenie liczy się w oknie,
 * ponieważ komenda `assistant.action.status` nie ma pola stanu w żądaniu
 * i wykaz zleceń przychodzi z rdzenia w całości.
 */
export interface FiltrZlecen {
  /** Kontrolka osadzana w nagłówku okna. */
  element: HTMLElement;
  /** Czy zlecenie przechodzi przez bieżące zawężenie. */
  przepusc(zlecenie: AssistantAction): boolean;
  /** Czy filtr cokolwiek zawęża — rozstrzyga zdanie pustki okna. */
  zawezony(): boolean;
  /** Nazwa bieżącego zawężenia; wchodzi do zdania pustki. */
  nazwa(): string;
}

/**
 * Stany kontraktu objęte każdą grupą filtra; grupa pusta znaczy brak
 * zawężenia. Anulowane stoją osobno od nieudanych, bo anulowanie jest decyzją
 * Operatora, a nie niepowodzeniem zlecenia.
 */
const GRUPY: Readonly<Record<KodFiltraZlecen, readonly AssistantActionStatus[]>> = {
  wszystkie: [],
  wToku: [
    AssistantActionStatus.Queued,
    AssistantActionStatus.Running,
    AssistantActionStatus.Paused,
  ],
  zakonczone: [AssistantActionStatus.Done],
  nieudane: [AssistantActionStatus.Failed],
  anulowane: [AssistantActionStatus.Cancelled],
};

export function utworzFiltrZlecen(naZmiane: () => void): FiltrZlecen {
  const kontrolka = wybor('Filtr statusu zleceń', POZYCJE_FILTRA);
  // Wykaz otwiera się bez zawężenia; zlecenie zakończone chwilę wcześniej jest tak samo istotne.
  kontrolka.value = 'wszystkie';
  kontrolka.addEventListener('change', naZmiane);

  const kod = (): KodFiltraZlecen => kontrolka.value as KodFiltraZlecen;

  return {
    element: kontrolka,
    przepusc(zlecenie) {
      const stany = GRUPY[kod()];
      return stany.length === 0 || stany.includes(zlecenie.status);
    },
    zawezony: () => GRUPY[kod()].length !== 0,
    nazwa: () => NAZWY_FILTRA[kod()],
  };
}
