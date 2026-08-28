import type { Agent } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { przyciskBezKomendy } from '../../modele/kontrolki-formularza-braki';

/**
 * Kontrolki i zdania wykazu archiwum ekspertów: wiersz pozycji oraz przekład
 * odmowy rdzenia na zdanie dla Operatora. Imię eksperta bierze się z nazwy
 * własnej, a wobec jej braku z nazwy technicznej.
 */
export function imie(ekspert: Agent | null): string {
  if (ekspert === null) return 'bez wybranego eksperta';
  return ekspert.displayName ?? ekspert.name;
}

/**
 * Zdanie odmowy złożone z trzech części: co odmówiło, dlaczego i co z tym
 * zrobić.
 *
 * Rdzeń bez podanej przyczyny też dostaje zdanie — cisza po naciśnięciu
 * przycisku jest gorsza niż odmowa bez szczegółu.
 */
export function odmowaRdzenia(czynnosc: string, kod?: string, tresc?: string): string {
  const powod = tresc ?? 'rdzeń nie podał przyczyny';
  const znak = kod === undefined ? '' : ` [${kod}]`;
  return `${czynnosc} — rdzeń odmówił${znak}: ${powod}. Sprawdź, czy ekspert nadal istnieje, i spróbuj ponownie.`;
}

/**
 * Wiersz wykazu z kontrolką pozycji.
 *
 * Kontrolka jest prawdziwa, gdy rdzeń zna komendę; w przeciwnym razie nazywa
 * brak, zamiast wyglądać na czynną.
 */
export function wierszWykazu(
  opis: string,
  etykieta: string,
  rdzenMa: boolean,
  komenda: string,
  czynnosc: () => Promise<void>,
): HTMLElement {
  const tresc = document.createElement('span');
  tresc.className = 'da-archiwum__opis';
  tresc.textContent = opis;

  const kontrolka = rdzenMa
    ? przycisk(etykieta, 'dn-btn dn-btn--zarys')
    : przyciskBezKomendy(etykieta, `Rdzeń nie wpina komendy ${komenda}.`);
  if (rdzenMa) kontrolka.addEventListener('click', () => void czynnosc());

  const element = document.createElement('li');
  element.className = 'da-archiwum__wiersz';
  element.append(tresc, kontrolka);
  return element;
}
