/**
 * Kolejność wczytania warstw wizualnych interfejsu. Kolejność importów jest
 * łańcuchem zależności: żetony motywu dają barwy, typografię, przestrzeń i ruch;
 * biblioteka `.dn-*` stoi na żetonach; `powloka.css` sprowadza zmienne widoku
 * okna do żetonów motywu, więc idzie po arkuszu tego widoku; `aplikacja.css`
 * osadza reguły w obszarze roboczym powłoki, więc jest ostatni.
 *
 * Arkusz kompletu sterowania (`sterowanie/sterowanie.css`) wciąga sam moduł
 * panelu — trafia do pakietu za tym wykazem, jako ostatni styl widoku.
 *
 * Plik jest wyłącznie wykazem importów wywołujących skutek uboczny; nie
 * eksportuje niczego i nie zawiera logiki.
 */

import '../motyw/motyw.css';
import '../komponenty/indeks.css';
import '../okno-komunikacji/okno.css';
import './powloka.css';
import './aplikacja.css';
