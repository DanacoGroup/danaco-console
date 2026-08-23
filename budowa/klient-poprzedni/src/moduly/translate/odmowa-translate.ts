import { ErrorCode } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';

/**
 * Odmowa czynności modelowej modułu Translate — powód rdzenia plus zdanie
 * o tym, czego brakuje.
 *
 * Dwie komendy modułu wykonuje rdzeń modelem: `translate.target.add` (przekłada
 * tekst źródłowy na język panelu) i `translate.source.detect` (pyta model
 * o język tekstu). Gdy rejestr nie ma ani jednego czynnego kanału, rdzeń odmawia
 * obu z kodem `channel_unavailable`. Sama wiadomość rdzenia jest powodem, ale nie
 * mówi Operatorowi, co ma zrobić, a jej brzmienie o niewpiętym rejestrze sugeruje
 * usterkę wewnętrzną tam, gdzie stoi zwykły brak włączonego kanału. Zdanie
 * o brakującej rzeczy dokłada klient — nie zastępując powodu rdzenia, tylko
 * idąc za nim.
 *
 * Zdania nie dokłada wspólny tłumacz kodów, bo ten sam kod `channel_unavailable`
 * niesie w tym module dwie różne odmowy: powyższą oraz odmowę
 * `translate.glossary.import`, gdzie rdzeń nie czyta plików z dysku Operatora.
 * Zdanie o kanale modelu byłoby przy imporcie nieprawdziwe, więc dokłada je
 * wyłącznie wywołanie, które wie, że jego czynność idzie modelem.
 */

/**
 * Czego brakuje, gdy rdzeń odmawia czynności modelowej z braku kanału.
 *
 * Zdanie mówi też, dlaczego Operator nie ma tu pola wyboru: żądania
 * `translate.target.add` i `translate.source.detect` nie niosą `channelId`,
 * więc rdzeń bierze kanał domyślny czynny i nie zgaduje niczego innego. Klient
 * nie wymyśla pól kontraktu, więc wyboru kanału tutaj nie dokłada.
 */
export const BRAK_KANALU_MODELU =
  'Brakuje czynnego kanału modelu. Przekład panelu i rozpoznanie języka rdzeń wykonuje ' +
  'modelem, a żądanie nie niesie wyboru kanału — kontrakt nie ma pola channelId dla ' +
  'translate.target.add ani translate.source.detect, więc rdzeń bierze kanał domyślny ' +
  'czynny. Włącz kanał w rejestrze kanałów modelu i powtórz czynność.';

/**
 * Zdanie odmowy czynności modelowej.
 *
 * Powód rdzenia idzie pierwszy i w całości — klient go nie skraca ani nie
 * podmienia. Zdanie o brakującym kanale dokłada się wyłącznie przy kodzie
 * `channel_unavailable`; przy każdym innym kodzie odmowa zostaje taka, jaką
 * podał rdzeń.
 */
export function zdanieOdmowyModelu(
  czynnosc: string,
  blad?: { code?: string; message?: string },
): string {
  const zdanie = opisOdmowy(czynnosc, blad?.code, blad?.message);
  if ((blad?.code ?? '') !== ErrorCode.ChannelUnavailable) return zdanie;
  return `${zdanie}. ${BRAK_KANALU_MODELU}`;
}
