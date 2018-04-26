package libs


var Prefix = "CAPTHCA_"

type RedisStore struct {

}

func (s *RedisStore) Set(id string, digits []byte) {
	SaveObject(Prefix+id,digits,60*10)
}

func (s *RedisStore) Get(id string, clear bool) (digits []byte) {
	if clear{
		b,_:= GetObjectAndCollect(Prefix+id)
		return b
	}else{
		b,_:= GetObjectAndDelay(Prefix+id,60*10)
		return b
	}

}

