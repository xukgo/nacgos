#ifndef __INNER_NACDEF_H__
#define __INNER_NACDEF_H__

typedef void (*ConfigUpdateEvent)(char *group, char *dataId, char *data);

typedef struct {
  char *name;
  ConfigUpdateEvent event;
} MatchVarEventHandler;

typedef struct {
  int count;
  MatchVarEventHandler* handlers;
} MatchVarEventHandlerCollection;

void nacdef_doUpdateEvent(ConfigUpdateEvent evt, char *group, char *dataId, char *data);

#endif