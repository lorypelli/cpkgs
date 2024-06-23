#ifndef INPUT_H
#define INPUT_H

#include <stdarg.h>
#include <stdbool.h>
#include <stddef.h>

// macro funzione file

#define f_in(f) \
  in(f, NULL)

// macro funzione stringa

#define s_in(s, ...) \
  in(NULL, s, ##__VA_ARGS__)

// macro funzione booleano

#define b_in(s, ...)                             \
  ({                                             \
    char b = c(s_in(s "(y/N) ", ##__VA_ARGS__)); \
    b == 'y' || b == 'Y';                        \
  })

// tipo stringa

typedef char *str;

// puntatore stringa

static str string = NULL;

// lunghezza stringa

static long length = 0;

// posizione carattere

static size_t position = 0;

// prototipi delle funzioni

static str input(const str, str, va_list);
extern bool v_str(const str, str);
extern str ln(const str);
extern size_t len(void);
extern void c_str(str, const size_t, const char);
extern void f_cat(const str, const size_t, const char, str, ...);
extern void f_cpy(const str, str);
extern void f_repeat(const str, const char, const size_t, const bool);
extern str f_replace(const str, const char, const char);
extern str in(const str, str, ...);
extern int sel_in(const size_t, str, str, ...);
extern char c(const str);
extern int i(const str);
extern long l(const str);
extern long long ll(const str);
extern float f(const str);
extern double d(const str);

#endif /* INPUT_H */