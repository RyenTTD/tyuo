mod database;

pub fn goodbye() {
    println!("Goodbye, world!");
    database::oh_no();
}

pub struct Model {
    database_manager: Box<database::DatabaseManager>,
    //dictionary_banned
    
    //contexts {id: {model(database), dictionary(database), dictionary_banned(ref)}}
}
impl Model {
    pub fn prepare(db_dir:&std::path::Path, banned_tokens_list:&std::path::Path) -> Model {
        //TODO: dev test
        println!("{}", database::DatabaseManager::prepare(db_dir).exists("hi"));
        
        return Model{
            database_manager: database::DatabaseManager::prepare(db_dir),
        };
    }
}
