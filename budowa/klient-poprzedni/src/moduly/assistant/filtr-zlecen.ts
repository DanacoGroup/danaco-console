import { AssistantActionStatus, type AssistantAction } from '../../../../shared/contract';
import { wybor } from '../../modele/kontrolki-formularza';
import { NAZWY_FILTRA, POZYCJE_FILTRA, type KodFiltraZlecen } from './etykiety-assistant';

/**
 * Filtr statusu wykazu zleceń Actions Monitora.
 *
 * Zawężenie liczy się w oknie, a nie w rdzeniu: `assistant.action.status` nie
 * ma pola stanu w żądaniu (`adapter_modul_asystent_czynnosci.go` czyta po
 * `windowId` albo po `actionId`), więc wykaz i tak przychodzi w całości.
 * Odczyt na każdą zmianę pozycji listy byłby wywołaniem bez nowej treści.
 *
 * Grupy filtra są złożone ze stanów kontraktu, nie z własnych nazw stanów.
 * „Anulowane" stoi osobno od „nieudanych", bo anulowanie jest decyzją
 * Operatora, a nie niepowodzeniem zlecenia — zlanie obu w jedną pozycję
 * kazałoby czytać własną decyzję jako usterkę.
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

/** Stany kontraktu objęte każdą grupą filtra; puste znaczy „bez zawężenia". */
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
  // Wykaz otwiera się bez zawężenia. Zlecenie zakończone chwilę wcześniej jest
  // dla Operatora wchodzącego do modułu tak samo istotne jak zlecenie w toku,
  // a wykaz zawężony od razu wyglądałby na pusty rdzeń.
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
