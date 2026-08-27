import { ErrorCode } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';

/**
 * Zdanie opisujące brak czynnego kanału modelu, dokładane do odmowy rdzenia,
 * gdy żądanie modelowe nie ma kanału do wyboru i rdzeń bierze kanał domyślny.
 */
export const BRAK_KANALU_MODELU =
  'Brakuje czynnego kanału modelu. Przekład panelu i rozpoznanie języka rdzeń wykonuje ' +
  'modelem, a żądanie nie niesie wyboru kanału — kontrakt nie ma pola channelId dla ' +
  'translate.target.add ani translate.source.detect, więc rdzeń bierze kanał domyślny ' +
  'czynny. Włącz kanał w rejestrze kanałów modelu i powtórz czynność.';

/**
 * Buduje zdanie odmowy czynności modelowej, zachowując powód rdzenia w całości
 * i dokładając zdanie o brakującym kanale wyłącznie przy odmowie z powodu
 * niedostępności kanału.
 */
export function zdanieOdmowyModelu(
  czynnosc: string,
  blad?: { code?: string; message?: string },
): string {
  const zdanie = opisOdmowy(czynnosc, blad?.code, blad?.message);
  if ((blad?.code ?? '') !== ErrorCode.ChannelUnavailable) return zdanie;
  return `${zdanie}. ${BRAK_KANALU_MODELU}`;
}
